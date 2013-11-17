using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Reflection;
using System.Reflection.Emit;
using System.Diagnostics;

namespace famousfront.utils
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
}
