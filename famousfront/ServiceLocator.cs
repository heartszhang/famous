using System.CodeDom;
using System.Linq;
using famousfront.core;
using famousfront.datamodels;
using famousfront.messages;
using famousfront.utils;
using famousfront.viewmodels;
using GalaSoft.MvvmLight.Messaging;
using Newtonsoft.Json;
using System;
using System.Diagnostics;
using System.IO;
using System.Threading.Tasks;

namespace famousfront
{
  internal static class BackendService
  {
    internal static string Compile(string address, string pattern, object paras = null, object paths = null)
    {
      var ps = new string[2]{pattern, null};

      if (paths != null)
      {
        ps[0] = paths.QueryPath(pattern);
      }
      if (paras != null)
      {
        ps[1] = paras.QueryString();
      }
      var rel = string.Join("?", ps.Where(s => !string.IsNullOrEmpty(s)));
      return "http://" + address + rel;
    }

    internal const string Tick = "/api/tick.json";
    internal const string FeedSourceSubscribe = "/api/feed_source/subscribe.json";
    internal const string FeedSourceFind = "/api/feed_source/find.json";
    internal const string ImageDimension = "/api/image/dimension.json";
    internal const string ImageDescription = "/api/image/description.json";
    internal const string Meta = "/api/meta.json";
    internal const string UpdatePopup = "/api/update/popup.json";
    internal const string FeedEntryUnread = "/api/feed_entry/unread.json";
    internal const string FeedEntrySourceUnreadcount = "/api/feed_entry/source/unread_count.json";
    internal const string FeedEntryFulldoc = "/api/feed_entry/fulldoc.json";
    internal const string FeedSourceAll = "/api/feed_source/all.json";
    internal const string FeedSourceShow = "/api/feed_source/show.json";
    internal const string FeedSourceUnsubscribe = "/api/feed_source/unsubscribe.json";
    internal const string ImageIcon = "/api/image/icon";

  }
  internal class ServiceLocator : CommandsService
  {
    static FeedsBackendConfig _backend;
    static FrontFlags _flags = new FrontFlags();
    public static FrontFlags FrontFlags
    {
      get { return _flags; }
    }

    private static MainViewModel _main;
    public static MainViewModel MainViewModel
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

    internal static void ClearMain()
    {
      if (_main != null)
        _main.Cleanup();
      _main = null;
    }

    static void CreateMain()
    {
      if (_main == null)
      {
        _main = new MainViewModel();
      }
    }

    internal static async void Startup()
    {
      CreateMain();
      Messenger.Default.Send(new famousfront.messages.BackendInitializing() { reason = Strings.Connecting });
      _flags = await Task.Run(() => DoLoad());
      await Task.Run(() =>
      {
        using (var p = Process.Start(BackendModule))
        {
          Messenger.Default.Send(new messages.BackendError { reason = BackendModule });
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
        var uri = BackendService.Compile(BackendAddress(), BackendService.Meta);

        var s = await HttpClientUtils.Get<FeedsBackendConfig>(uri);
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
      ShutdownBackend();
      Messenger.Default.Send(new BackendShutdown() { });
    }

    private static void DoSave()
    {
      var js = new JsonSerializer();
      using (var writer = new StreamWriter(ConfigFile))
      using (var jwriter = new JsonTextWriter(writer))
      {
        js.Serialize(writer, FrontFlags);
      }
    }
    private static FrontFlags DoLoad()
    {
      FrontFlags v = null;
      var c = ConfigFile;
      if (!File.Exists(c))
        return new FrontFlags();
      var js = new JsonSerializer();
      using (var reader = new StreamReader(c))
      using (var jreader = new JsonTextReader(reader))
      {
        v = js.Deserialize<FrontFlags>(jreader);
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
    
    internal static string BackendAddress()
    {
      return _flags.Backend;
      //return BackendService.Compile(_flags.Backend, rel);
      //return "http://" + _flags.Backend + rel;
    }
   
    internal static void Log(string format, params object[] args){
      var m = string.Format(format, args);
      Messenger.Default.Send(new messages.BackendError { reason = m });
    }
  }
}
