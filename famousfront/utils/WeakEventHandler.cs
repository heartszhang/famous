using System;

namespace famousfront.utils
{
  internal delegate void UnregisterCallback<E>(EventHandler<E> eventHandler)
    where E : EventArgs;

  internal interface IWeakEventHandler<TE>
    where TE : EventArgs
  {
    EventHandler<TE> Handler { get; }
  }

  internal class WeakEventHandler<T, TE> : IWeakEventHandler<TE>
    where T : class
    where TE : EventArgs
  {
    private delegate void OpenEventHandler(T @this, object sender, TE e);

    private readonly WeakReference m_target_ref;
    private readonly OpenEventHandler m_open_handler;
    private readonly EventHandler<TE> m_handler;
    private UnregisterCallback<TE> m_unregister;

    public WeakEventHandler(EventHandler<TE> eventHandler, UnregisterCallback<TE> unregister)
    {
      m_target_ref = new WeakReference(eventHandler.Target);
      m_open_handler = (OpenEventHandler)Delegate.CreateDelegate(typeof(OpenEventHandler),
        null, eventHandler.Method);
      m_handler = Invoke;
      m_unregister = unregister;
    }

    public void Invoke(object sender, TE e)
    {
      var target = (T)m_target_ref.Target;

      if (target != null)
        m_open_handler.Invoke(target, sender, e);
      else if (m_unregister != null)
      {
        m_unregister(m_handler);
        m_unregister = null;
      }
    }

    public EventHandler<TE> Handler
    {
      get { return m_handler; }
    }

    public static implicit operator EventHandler<TE>(WeakEventHandler<T, TE> weh)
    {
      return weh.m_handler;
    }
  }
  internal interface IWeakEventHandlerSpecial<out TEventHandler>
  {
    TEventHandler Handler { get; }
  }
}
