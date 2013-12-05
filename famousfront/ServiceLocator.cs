using famousfront.core;
using famousfront.datamodels;
using famousfront.messages;
using famousfront.utils;
using famousfront.viewmodels;
using GalaSoft.MvvmLight.Command;
using GalaSoft.MvvmLight.Messaging;
//using MetroLog;
using Newtonsoft.Json;
using System;
using System.Diagnostics;
using System.Diagnostics.CodeAnalysis;
using System.IO;
using System.Threading.Tasks;
using System.Windows.Input;

namespace famousfront
{
  internal class ServiceLocator
  {
    //        private static ILogger _log = LogManagerFactory.DefaultLogManager.GetLogger<ServiceLocator>();
    private static MainViewModel _main;

    private static SettingsViewModel _settings;
    static Flags _flags = new Flags();
    static FeedsBackendConfig _backend;

    public static Flags Flags
    {
      get { return _flags; }
    }
    /// <summary>
    /// Gets the Main property.
    /// </summary>
    internal static MainViewModel Main
    {
      get
      {
        if (_main == null)
        {
          CreateMain();
        }

        return _main;
      }
    }

    /// <summary>
    /// Gets the Settings property.
    /// </summary>
    internal static SettingsViewModel Settings
    {
      get
      {
        if (_settings == null)
        {
          CreateSettings();
        }

        return _settings;
      }
    }

    /// <summary>
    /// Gets the Main property.
    /// </summary>
    [SuppressMessage("Microsoft.Performance",
        "CA1822:MarkMembersAsStatic",
        Justification = "This non-static member is needed for data binding purposes.")]
    internal MainViewModel MainViewModel
    {
      get
      {
        return Main;
      }
    }

    /// <summary>
    /// Gets the Settings property.
    /// </summary>
    [SuppressMessage("Microsoft.Performance",
        "CA1822:MarkMembersAsStatic",
        Justification = "This non-static member is needed for data binding purposes.")]
    internal SettingsViewModel SettingsViewModel
    {
      get
      {
        return Settings;
      }
    }

    /// <summary>
    /// Provides a deterministic way to delete the Main property.
    /// </summary>
    internal static void ClearMain()
    {
      if (_main != null)
        _main.Cleanup();
      _main = null;
    }

    /// <summary>
    /// Provides a deterministic way to delete the Settings property.
    /// </summary>
    internal static void ClearSettings()
    {
      if (_settings != null)
        _settings.Cleanup();
      _settings = null;
    }

    /// <summary>
    /// Provides a deterministic way to create the Main property.
    /// </summary>
    internal static void CreateMain()
    {
      if (_main == null)
      {
        _main = new MainViewModel();
      }
    }

    /// <summary>
    /// Provides a deterministic way to create the Settings property.
    /// </summary>
    internal static void CreateSettings()
    {
      if (_settings == null)
      {
        _settings = new SettingsViewModel();
      }
    }
    internal static async void Startup()
    {
      CreateMain();
      Messenger.Default.Send(new famousfront.messages.BackendInitializing() { reason = Strings.Connecting });
      var v = await Task.Run(() => DoLoad());
      if (v != null) _flags = v;
      await Task.Run(() =>
      {
        using (var p = Process.Start(BackendModule))
        {
        }
      });
      await StartKeepaliver();
      Messenger.Default.Send(new famousfront.messages.BackendInitialized() { reason = Strings.OK });
    }
    static async Task StartKeepaliver()
    {
      await DoKeepalive();
    }
    static async Task DoKeepalive()
    {
      for (; ; )
      {
        var s = await HttpClientUtils.Get<FeedsBackendConfig>(BackendPath("/api/meta.json"));
        if (s.code == 0)
        {
          _backend = s.data;
          break;
        }
        Messenger.Default.Send(new BackendInitializing() { reason = s.reason });
        await Task.Delay(_flags.KaPeriod);
      }
    }
    static void ShutdownBackend()
    {
    }
    internal static void Shutdown()
    {
      ClearMain();
      ClearSettings();
      ShutdownBackend();
      Messenger.Default.Send(new BackendShutdown() { });
    }

    private static void DoSave()
    {
      var js = new JsonSerializer();
      using (var writer = new StreamWriter(ConfigFile))
      using (var jwriter = new JsonTextWriter(writer))
      {
        js.Serialize(writer, Flags);
      }
    }
    private static Flags DoLoad()
    {
      Flags v = null;
      var c = ConfigFile;
      if (!File.Exists(c))
        return v;
      var js = new JsonSerializer();
      using (var reader = new StreamReader(c))
      using (var jreader = new JsonTextReader(reader))
      {
        v = js.Deserialize<Flags>(jreader);
      }
      return v;
    }
    public static string ConfigFolder
    {
      get { return Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData), "famous"); }
    }
    public static string ConfigFile
    {
      get { return Path.Combine(ConfigFolder, "famous.config"); }
    }
    public static string RootFolder
    {
      get { return AppDomain.CurrentDomain.BaseDirectory; }
    }
    internal static string BackendModule
    {
      get { return Path.Combine(RootFolder, "backend.bat"); }
    }
    internal static string BackendPath(string rel)
    {
      return "http://" + _flags.Backend + rel;
    }

    ICommand _toggle_feedsources_view;
    public ICommand ToggleFeedSourcesViewCommand
    {
      get { return _toggle_feedsources_view ?? (_toggle_feedsources_view = toggle_feedsources_view()); }
    }
    ICommand toggle_feedsources_view()
    {
      return new RelayCommand(ExecuteToggleFeedSource);
    }
    void ExecuteToggleFeedSource()
    {
      Messenger.Default.Send(new ToggleFeedSource());
    }

    ICommand _show_find_feedsource_view;
    public ICommand FindFeedSourceCommand
    {
      get { return _show_find_feedsource_view ?? (_show_find_feedsource_view = show_find_feedsource_view()); }
    }
    ICommand show_find_feedsource_view()
    {
      return new RelayCommand(ExecuteShowFindFeedSourceView);
    }
    void ExecuteShowFindFeedSourceView()
    {
      Messenger.Default.Send(new ShowFindFeedSourceView());
    }

    ICommand _hyperlink_navigate;
    public ICommand HyperlinkNavigateCommand
    {
      get { return _hyperlink_navigate ?? (_hyperlink_navigate = hyperlink_navigate()); }
    }
    ICommand hyperlink_navigate()
    {
      return new RelayCommand<Uri>(ExecuteHyperlinkNavigate);
    }
    void ExecuteHyperlinkNavigate(Uri url)
    {
      if (url == null)
        return;
      using (var p = Process.Start(url.ToString())) { };
    }
  }
}
