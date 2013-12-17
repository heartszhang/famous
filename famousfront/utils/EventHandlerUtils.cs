using System;
using System.Diagnostics;
using famousfront.Properties;

namespace famousfront.utils
{
  internal static class EventHandlerUtils
  {
    public static EventHandler<TE> MakeWeak<TE>(this EventHandler<TE> eventHandler, UnregisterCallback<TE> unregister)
      where TE : EventArgs
    {
      if (eventHandler == null)
        throw new ArgumentNullException("eventHandler");
      if (eventHandler.Method.IsStatic || eventHandler.Target == null)
        throw new ArgumentException(Resources.EventHandlerUtils_MakeWeak_Only_instance_methods_are_supported_, "eventHandler");

      var wehType = typeof(WeakEventHandler<,>).MakeGenericType(eventHandler.Method.DeclaringType, typeof(TE));
      var wehConstructor = wehType.GetConstructor(new[] { typeof(EventHandler<TE>), 
        typeof(UnregisterCallback<TE>) });

      Debug.Assert(wehConstructor != null);
      var weh = (IWeakEventHandler<TE>)wehConstructor.Invoke(
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
      var wehType = typeof(WeakEventHandlerSpecial<,,>).MakeGenericType(ehDelegate.Method.DeclaringType, typeof(TEventHandler), eventArgsType);

      var wehConstructor = wehType.GetConstructor(new[] { typeof(Delegate), typeof(Action<object>) });

      Debug.Assert(wehConstructor != null, "Something went wrong. There should be constructor with these types");

      var weh = (IWeakEventHandlerSpecial<TEventHandler>)wehConstructor.Invoke(new object[] { eventHandler, (Action<object>)(o => unregister((TEventHandler)o)) });

      return weh.Handler;
    }
  }
}