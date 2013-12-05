using famousfront.core;
using famousfront.messages;

namespace famousfront.viewmodels
{
  internal enum MainViewStatus
  {
    OFFLINE,
    READY,
    BOOTSTRAP,
    CONNECTING,
  }
  internal class MainViewModel : TaskViewModel
  {
    internal MainViewModel()
    {
      MessengerInstance.Register<BackendInitialized>(this, OnBackendInitialized);
      MessengerInstance.Register<ShowFindFeedSourceView>(this, ExecuteShowFindFeedSourceView);
      _content = new ContentViewModel();
      _offline = new OfflineViewModel();
      _booting = new BootstrapViewModel();
    }
    public override void Cleanup()
    {
      _content.Cleanup();
      _offline.Cleanup();
      _booting.Cleanup();
      base.Cleanup();
    }

    ContentViewModel _content;
    public ContentViewModel ContentViewModel
    {
      get { return _content; }
      internal set
      {
        Set(ref _content, value);
      }
    }

    BootstrapViewModel _booting;
    public BootstrapViewModel BootstrapViewModel
    {
      get { return _booting; }
      internal set
      {
        Set(ref _booting, value);
      }
    }

    OfflineViewModel _offline;
    public OfflineViewModel OfflineViewModel
    {
      get { return _offline; }
      internal set
      {
        Set(ref _offline, value);
      }
    }
    SettingsViewModel _settings;
    public SettingsViewModel SettingsViewModel
    {
      get { return _settings; }
      internal set
      {
        Set(ref _settings, value);
      }
    }

    FeedSourceFindViewModel _feedsourcefind;
    public FeedSourceFindViewModel FeedSourceFindViewModel
    {
      get { return _feedsourcefind; }
      internal set { Set(ref _feedsourcefind, value); }
    }
    bool _show_settings_view = false;
    public bool ShowSettingsView
    {
      get { return _show_settings_view; }
      internal set
      {
        Set(ref _show_settings_view, value);
      }
    }
    bool _show_feedsourcefind_view;
    public bool ShowFeedSourceFindView
    {
      get { return _show_feedsourcefind_view; }
      internal set { Set(ref _show_feedsourcefind_view, value); }
    }
    MainViewStatus _status = MainViewStatus.BOOTSTRAP;
    public MainViewStatus Status
    {
      get { return _status; }
      internal set
      {
        Set(ref _status, value);
      }
    }

    void OnBackendInitialized(BackendInitialized msg)
    {
      Status = MainViewStatus.READY;
      _content.ReloadFeedSources();
    }

    void ExecuteShowFindFeedSourceView(ShowFindFeedSourceView msg)
    {
      if (FeedSourceFindViewModel == null)
      {
        FeedSourceFindViewModel = new FeedSourceFindViewModel();
      }

      ShowFeedSourceFindView = !ShowFeedSourceFindView;
      if (!ShowFeedSourceFindView)
      {
        //            FeedSourceFindViewModel = null;  // reserve previous search
      }
    }
  }
}
