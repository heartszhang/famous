using famousfront.messages;

namespace famousfront.viewmodels
{
  class BootstrapViewModel : famousfront.core.TaskViewModel
  {
    internal BootstrapViewModel()
    {
      MessengerInstance.Register<BackendInitializing>(this, OnBackendInitializing);
      MessengerInstance.Register<BackendInitialized>(this, OnBackendInitialized);
    }
    void OnBackendInitialized(BackendInitialized msg)
    {
      IsBusying = false;
      Reason = msg.reason;
    }
    void OnBackendInitializing(BackendInitializing msg)
    {
      IsBusying = true;
      Reason = msg.reason;
    }
  }
}
