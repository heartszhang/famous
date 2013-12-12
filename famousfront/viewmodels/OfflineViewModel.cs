using GalaSoft.MvvmLight.Threading;
namespace famousfront.viewmodels
{
  using Messages = System.Collections.ObjectModel.ObservableCollection<string>;
  class OfflineViewModel : famousfront.core.DialogViewModel
  {
    internal OfflineViewModel()
      : base(CommandFlags.HasNone)
    {

    }
  }
  class MessagesViewModel : famousfront.core.OverlappedViewModel
  {
    readonly Messages _messages = new Messages();
    public System.ComponentModel.ICollectionView Messages
    {
      get { return System.Windows.Data.CollectionViewSource.GetDefaultView(_messages); }
    }
    internal MessagesViewModel()
    {
      MessengerInstance.Register<famousfront.messages.BackendError>(this, OnBackendError);
    }
    protected override void ExecuteClose()
    {
      base.ExecuteClose();
      MessengerInstance.Send(new famousfront.messages.ShowMessagesView());
    }
    private void OnBackendError(messages.BackendError obj)
    {
      DispatcherHelper.CheckBeginInvokeOnUI(() =>
      {
        _messages.Add(obj.reason);
        if (_messages.Count > 1000)
        {
          while (_messages.Count > 500)
          {
            _messages.RemoveAt(0);
          }
        }
      });
    }
  }
  }
