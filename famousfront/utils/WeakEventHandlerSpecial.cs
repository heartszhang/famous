using System;

namespace famousfront.utils
{
  internal class WeakEventHandlerSpecial<T, TEventHandler, TEventArgs> : IWeakEventHandlerSpecial<TEventHandler>
    where T : class
    where TEventArgs : EventArgs
  {
    private delegate void OpenEventHandler(T @this, object sender, TEventArgs e);

    private readonly WeakReference _targetRef;
    private readonly OpenEventHandler _openHandler;
    private Action<object> _unregister;
    private readonly TEventHandler _handler;

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
}