﻿using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Reflection;
using System.Reflection.Emit;
using System.Diagnostics;

namespace famousfront.core
{
    internal delegate void UnregisterCallback<E>(EventHandler<E> eventHandler)
      where E : EventArgs;

    internal interface IWeakEventHandler<E>
      where E : EventArgs
    {
        EventHandler<E> Handler { get; }
    }

    internal class WeakEventHandler<T, E> : IWeakEventHandler<E>
        where T : class
        where E : EventArgs
    {
        private delegate void OpenEventHandler(T @this, object sender, E e);

        private WeakReference m_TargetRef;
        private OpenEventHandler m_OpenHandler;
        private EventHandler<E> m_Handler;
        private UnregisterCallback<E> m_Unregister;

        public WeakEventHandler(EventHandler<E> eventHandler, UnregisterCallback<E> unregister)
        {
            m_TargetRef = new WeakReference(eventHandler.Target);
            m_OpenHandler = (OpenEventHandler)Delegate.CreateDelegate(typeof(OpenEventHandler),
              null, eventHandler.Method);
            m_Handler = Invoke;
            m_Unregister = unregister;
        }

        public void Invoke(object sender, E e)
        {
            T target = (T)m_TargetRef.Target;

            if (target != null)
                m_OpenHandler.Invoke(target, sender, e);
            else if (m_Unregister != null)
            {
                m_Unregister(m_Handler);
                m_Unregister = null;
            }
        }

        public EventHandler<E> Handler
        {
            get { return m_Handler; }
        }

        public static implicit operator EventHandler<E>(WeakEventHandler<T, E> weh)
        {
            return weh.m_Handler;
        }
    }
    internal interface IWeakEventHandlerSpecial<TEventHandler>
    {
        TEventHandler Handler { get; }
    }

    internal class WeakEventHandlerSpecial<T, TEventHandler, TEventArgs> : IWeakEventHandlerSpecial<TEventHandler>
        where T : class
        where TEventArgs : EventArgs
    {
        private delegate void OpenEventHandler(T @this, object sender, TEventArgs e);

        private readonly WeakReference _targetRef;
        private readonly OpenEventHandler _openHandler;
        private Action<object> _unregister;
        private TEventHandler _handler;

        public WeakEventHandlerSpecial(Delegate eventHandler, Action<object> unregister)
        {
            _targetRef = new WeakReference(eventHandler.Target);
            _openHandler = (OpenEventHandler)Delegate.CreateDelegate(typeof(OpenEventHandler), null, eventHandler.Method);
            _unregister = unregister;

            var t = GetType();
            var mi = t.GetMethod("Invoke");
            _handler = (TEventHandler)((object)Delegate.CreateDelegate(typeof(TEventHandler), this, mi));
        }

        public void Invoke(object sender, TEventArgs e)
        {
            var target = (T)_targetRef.Target;

            if (target != null)
            {
                _openHandler.Invoke(target, sender, e);
            }
            else if (_unregister != null)
            {
                _unregister(_handler);
                _unregister = null;
            }
        }

        public TEventHandler Handler
        {
            get { return _handler; }
        }
    }

    internal static class EventHandlerUtils
    {
        public static EventHandler<E> MakeWeak<E>(this EventHandler<E> eventHandler, UnregisterCallback<E> unregister)
          where E : EventArgs
        {
            if (eventHandler == null)
                throw new ArgumentNullException("eventHandler");
            if (eventHandler.Method.IsStatic || eventHandler.Target == null)
                throw new ArgumentException("Only instance methods are supported.", "eventHandler");

            Type wehType = typeof(WeakEventHandler<,>).MakeGenericType(eventHandler.Method.DeclaringType, typeof(E));
            ConstructorInfo wehConstructor = wehType.GetConstructor(new Type[] { typeof(EventHandler<E>), 
      typeof(UnregisterCallback<E>) });

            IWeakEventHandler<E> weh = (IWeakEventHandler<E>)wehConstructor.Invoke(
              new object[] { eventHandler, unregister });

            return weh.Handler;
        }
        /// <summary>
        /// Analogous to MakeWeak method, but allows working with non-standard event handlers (i.e. that differs from EventHandler{TEventArgs}).
        /// Somewhat ugly code, but it works.
        /// </summary>
        public static TEventHandler MakeWeakSpecial<TEventHandler>(this TEventHandler eventHandler, Action<TEventHandler> unregister)
        {
            if (eventHandler == null)
                throw new ArgumentNullException("eventHandler");

            var ehDelegate = (Delegate)(object)eventHandler;
            var eventArgsType = ehDelegate.Method.GetParameters()[1].ParameterType;
            Type wehType = typeof(WeakEventHandlerSpecial<,,>).MakeGenericType(ehDelegate.Method.DeclaringType, typeof(TEventHandler), eventArgsType);

            ConstructorInfo wehConstructor = wehType.GetConstructor(new[] { typeof(Delegate), typeof(Action<object>) });

            Debug.Assert(wehConstructor != null, "Something went wrong. There should be constructor with these types");

            var weh = (IWeakEventHandlerSpecial<TEventHandler>)wehConstructor.Invoke(new object[] { eventHandler, (Action<object>)(o => unregister((TEventHandler)o)) });

            return weh.Handler;
        }
    }

    /*
    /// <summary>
        /// Weak delegates are useful where you don't want the lifetime of a delegate target
        /// to be bound to the delegate's lifetime.
        ///
        /// The WeakDelegateFactory is used for the the same purpose as WeakEventHandler, but
        /// can work for any delegate type. Is uses Reflection.Emit to create an exact replica
        /// of the WeakEventHandler class, but for the specific delegate type.
        /// </summary>
        public static class WeakDelegateFactory
        {
            private static AssemblyBuilder theDelegatesAsm;
            private static ModuleBuilder theModuleBuilder;

            public static TDelegate MakeWeak<TDelegate>(this TDelegate targetDelegate, Action<TDelegate> unregisterDelegate) where TDelegate : class
            {
                EnsureBuildersCreated();

                var type = targetDelegate.GetType();

                if (type.GetMethod("Invoke").ReturnType != typeof(void))
                    throw new ArgumentException("Weak delegates can only be created for delegates with void return type.", "targetDelegate");

                // create the class builder that will give birth to this weak delegate class
                var className = "WeakDelg->" + GetPrettyName(type);
                var classType = theModuleBuilder.GetType(className)
                                ?? CreateWeakDelegateClass(targetDelegate, type, className);

                var ctor = classType.GetConstructors()[0];
                var weakDelg = ctor.Invoke(new object[] { targetDelegate, unregisterDelegate });
                return (TDelegate)(object)Delegate.CreateDelegate(type, weakDelg, classType.GetMethod("Invoke"));
            }

            private static Type CreateWeakDelegateClass<TDelegate>(TDelegate targetDelegate, Type type, string className) where TDelegate : class
            {
                var classBuilder = theModuleBuilder.DefineType(className, TypeAttributes.Class | TypeAttributes.Public);

                // create the class fields for the unbound delegate, the weak reference to the delegate's target and an unregister callback
                var weakRefFld = classBuilder.DefineField("myWeakRef", typeof(WeakReference), FieldAttributes.Private);

                MethodBuilder unboundDelgInvoke;
                var unboundDelgType = CreateUnboundDelegateType(classBuilder, targetDelegate, out unboundDelgInvoke);
                var unboundDelgFld = classBuilder.DefineField("myUnboundDelg", unboundDelgType, FieldAttributes.Private);

                var unregisterDelgFld = classBuilder.DefineField("myUnregisterDelg", typeof(Action<TDelegate>), FieldAttributes.Private);

                // create constructor
                DefineConstructor<TDelegate>(type, classBuilder, unboundDelgType, weakRefFld, unboundDelgFld, unregisterDelgFld);

                // create the method that implements the weak delegate calling
                DefineInvokeMethod(targetDelegate, classBuilder, weakRefFld, unboundDelgFld, unregisterDelgFld, unboundDelgInvoke);

                unboundDelgType.CreateType();
                return classBuilder.CreateType();
            }

            private static void DefineConstructor<TDelegate>(Type delegateType, TypeBuilder classBuilder, TypeBuilder unboundDelgType, FieldBuilder weakRefFld, FieldBuilder unboundDelgFld, FieldBuilder unregisterDelgFld)
            {
                var target = delegateType.GetProperty("Target");
                // create the constructor; it initializes the unbound delegate and the weak reference
                var constructor = classBuilder.DefineConstructor(MethodAttributes.Public, CallingConventions.HasThis,
                                                                 new[] { delegateType, typeof(Action<TDelegate>) });
                var ctorIl = constructor.GetILGenerator();

                // call object.ctor()
                ctorIl.Emit(OpCodes.Ldarg_0);
                ctorIl.Emit(OpCodes.Call, typeof(object).GetConstructor(new Type[] { }));

                // store a WeakReference to the delegate's target into a field
                ctorIl.Emit(OpCodes.Ldarg_0);
                ctorIl.Emit(OpCodes.Ldarg_1);
                ctorIl.Emit(OpCodes.Callvirt, target.GetGetMethod());
                ctorIl.Emit(OpCodes.Newobj, typeof(WeakReference).GetConstructor(new[] { typeof(object) }));
                ctorIl.Emit(OpCodes.Stfld, weakRefFld);

                // create an unbound delegate type from the given delegate's method
                ctorIl.Emit(OpCodes.Ldarg_0); // for stfld
                ctorIl.Emit(OpCodes.Ldtoken, unboundDelgType);
                ctorIl.Emit(OpCodes.Call, typeof(Type).GetMethod("GetTypeFromHandle", new[] { typeof(RuntimeTypeHandle) })); // first parameter; load the unbound delegate's type
                ctorIl.Emit(OpCodes.Ldnull); // second parameter; no instance given, i.e. treat it as a static method
                ctorIl.Emit(OpCodes.Ldarg_1); // third parameter to CreateDelegate
                ctorIl.Emit(OpCodes.Callvirt, typeof(Delegate).GetProperty("Method").GetGetMethod());
                ctorIl.Emit(OpCodes.Call, typeof(Delegate).GetMethod("CreateDelegate", new[] { typeof(Type), typeof(object), typeof(MethodInfo) }));
                ctorIl.Emit(OpCodes.Castclass, unboundDelgType);
                ctorIl.Emit(OpCodes.Stfld, unboundDelgFld);

                // store the unregister callback
                ctorIl.Emit(OpCodes.Ldarg_0);
                ctorIl.Emit(OpCodes.Ldarg_2);
                ctorIl.Emit(OpCodes.Stfld, unregisterDelgFld);

                ctorIl.Emit(OpCodes.Ret);
            }

            private static void DefineInvokeMethod<TDelegate>(TDelegate targetDelegate, TypeBuilder classBuilder, FieldBuilder weakRefFld, FieldBuilder unboundDelgFld, FieldBuilder unregisterDelgFld, MethodBuilder unboundDelgInvoke) where TDelegate : class
            {
                var targetType = ((Delegate)(object)targetDelegate).Target.GetType();
                var unregisterDelgType = typeof(Action<TDelegate>);
                var delegateType = targetDelegate.GetType();
                var delegateSignature = delegateType.GetMethod("Invoke");
                var delgParamTypes = GetParameterTypes(delegateSignature);
                var invoker = classBuilder.DefineMethod("Invoke", MethodAttributes.Public, delegateSignature.ReturnType, delgParamTypes.ToArray());
                var invokerIl = invoker.GetILGenerator();

                var targetLocal = invokerIl.DeclareLocal(targetType);

                var endLabel = invokerIl.DefineLabel();
                var targetIsNullLabel = invokerIl.DefineLabel();

                // get Target from weak reference
                invokerIl.Emit(OpCodes.Ldarg_0);
                invokerIl.Emit(OpCodes.Ldfld, weakRefFld);
                invokerIl.Emit(OpCodes.Callvirt, typeof(WeakReference).GetProperty("Target").GetGetMethod());
                invokerIl.Emit(OpCodes.Isinst, targetType);
                invokerIl.Emit(OpCodes.Castclass, targetType);
                invokerIl.Emit(OpCodes.Stloc, targetLocal);
                invokerIl.Emit(OpCodes.Ldloc, targetLocal);
                invokerIl.Emit(OpCodes.Brfalse_S, targetIsNullLabel);

                // the target is not null - call the unbound delegate
                invokerIl.Emit(OpCodes.Ldarg_0);
                invokerIl.Emit(OpCodes.Ldfld, unboundDelgFld);
                invokerIl.Emit(OpCodes.Ldloc, targetLocal);
                invokerIl.Emit(OpCodes.Ldarg_1);
                invokerIl.Emit(OpCodes.Ldarg_2);
                invokerIl.Emit(OpCodes.Callvirt, unboundDelgInvoke);
                invokerIl.Emit(OpCodes.Ret);

                invokerIl.MarkLabel(targetIsNullLabel);
                // target was null, call unregister
                invokerIl.Emit(OpCodes.Ldarg_0);
                invokerIl.Emit(OpCodes.Ldfld, unregisterDelgFld);
                invokerIl.Emit(OpCodes.Brfalse_S, endLabel); // check if unregister is null
                invokerIl.Emit(OpCodes.Ldarg_0);
                invokerIl.Emit(OpCodes.Ldfld, unregisterDelgFld);
                invokerIl.Emit(OpCodes.Ldarg_0);
                invokerIl.Emit(OpCodes.Ldftn, invoker);
                invokerIl.Emit(OpCodes.Newobj, delegateType.GetConstructor(new[] { typeof(object), typeof(IntPtr) })); // create delegate from invoke
                invokerIl.Emit(OpCodes.Callvirt, unregisterDelgType.GetMethod("Invoke")); // call unregister
                invokerIl.Emit(OpCodes.Ldarg_0);
                invokerIl.Emit(OpCodes.Ldnull);
                invokerIl.Emit(OpCodes.Stfld, unregisterDelgFld); // nullify unregister
                invokerIl.MarkLabel(endLabel);

                invokerIl.Emit(OpCodes.Ret);
            }

            private static void EnsureBuildersCreated()
            {
                if (theDelegatesAsm == null)
                {
                    theDelegatesAsm = AppDomain.CurrentDomain.DefineDynamicAssembly(new AssemblyName("WeakDelegates"), AssemblyBuilderAccess.Run);
                    theModuleBuilder = theDelegatesAsm.DefineDynamicModule("WeakDelegates");
                }
            }

            /// <summary>
            /// Creates an unbound delegate type from the type of a normal delegate
            /// </summary>
            private static TypeBuilder CreateUnboundDelegateType(TypeBuilder parentClass, object boundDelegate, out MethodBuilder outInvokeMethod)
            {
                var invokeMethod = boundDelegate.GetType().GetMethod("Invoke");
                var targetType = ((Delegate)boundDelegate).Target.GetType();
                return CreateDelegateType(parentClass, invokeMethod.ReturnType, new[] { targetType }.Concat(GetParameterTypes(invokeMethod)).ToArray(), out outInvokeMethod);
            }

            private static IEnumerable<Type> GetParameterTypes(MethodInfo mi)
            {
                return from p in mi.GetParameters() select p.ParameterType;
            }

            /// <summary>
            /// Basically does what the compiler does when it has to compile a delegate statement.
            /// Creates the delegate as a nested type.
            /// </summary>
            /// <returns>A class derived from System.MulticastDelegate that can be called with the given parameters</returns>
            private static TypeBuilder CreateDelegateType(TypeBuilder parentClass, Type returnType, Type[] parameters, out MethodBuilder outInvokeMethod)
            {
                // from Joel Pobar's CLR weblog
                // Creating delegate types via Reflection.Emit
                // http://blogs.msdn.com/joelpob/archive/2004/02/15/73239.aspx

                var delgTypeBuilder = parentClass.DefineNestedType("UnboundDelegate", TypeAttributes.Class | TypeAttributes.NestedPublic | TypeAttributes.Sealed | TypeAttributes.AnsiClass | TypeAttributes.AutoClass, typeof(MulticastDelegate));
                var constructorBuilder = delgTypeBuilder.DefineConstructor(MethodAttributes.RTSpecialName | MethodAttributes.HideBySig | MethodAttributes.Public, CallingConventions.Standard, new[] { typeof(object), typeof(IntPtr) });
                constructorBuilder.SetImplementationFlags(MethodImplAttributes.Runtime | MethodImplAttributes.Managed);

                outInvokeMethod = delgTypeBuilder.DefineMethod("Invoke", MethodAttributes.Public | MethodAttributes.HideBySig | MethodAttributes.NewSlot | MethodAttributes.Virtual, returnType, parameters);
                outInvokeMethod.SetImplementationFlags(MethodImplAttributes.Runtime | MethodImplAttributes.Managed);

                return delgTypeBuilder;
            }

            private static string GetPrettyName(Type t)
            {
                if (t.IsGenericType)
                {
                    return t.Namespace + '.' + t.Name + '(' + String.Join(";", (from p in t.GetGenericArguments() select p.Name).ToArray()) + ')';
                }
                else
                {
                    return t.FullName;
                }
            }
        }
     */
}