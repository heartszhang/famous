using System.Collections.ObjectModel;
using GalaSoft.MvvmLight.Threading;

namespace famousfront.viewmodels
{
  class MessagesViewModel : famousfront.core.OverlappedViewModel
  {
    readonly ObservableCollection<string> _messages = new ObservableCollection<string>();
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