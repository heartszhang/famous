using System;
using System.Linq;
using famousfront.datamodels;
using System.Windows.Input;
using GalaSoft.MvvmLight.Command;
using famousfront.utils;
using famousfront.messages;
using System.Diagnostics;
namespace famousfront.viewmodels
{
  internal class FeedSourceViewModel : core.TaskViewModel
  {
    readonly ICommand _drop_self;
    readonly ICommand _goto_page;
    int _page;
    static readonly DateTime Utime = new DateTime(1970, 1, 1, 0, 0, 0, 0);
    readonly FeedSource _;
    internal FeedSourceViewModel(FeedSource val)
    {
      _ = val;
      News = _.description;
      _goto_page = new RelayCommand<int>(ExecuteGotoPage);
      _drop_self = new RelayCommand(ExecuteDropSelf);
      LoadUnreadCount();
    }
    public string Category { get { return _category ?? first_or_default_category(); } }
    public string Name { get { return _.name; } set { _.name = value; RaisePropertyChanged(); } }
    public string Uri { get { return _.uri; } }
    public string Description { get { return _.description; } }

    public ICommand GotoPageCommand
    {
      get { return _goto_page ; }
    }

    public ICommand DropSelfCommand
    {
      get { return _drop_self ; }
    }
    public string Logo { get { return logo(); }  }
    string logo()
    {
      var uri = _.logo;
      if (string.IsNullOrEmpty(uri))
        uri = _.website;
      //var rel = "/api/image/icon?uri=" + System.Uri.UnescapeDataString(hint);
      var x = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.ImageIcon, new{uri });
      return x;
    }
    string _news;
    public string News { get { return _news; } private set { Set(ref _news, value); } }

    public int UnreadCount
    {
      get { return _.unreaded; }
      private set { if (_.unreaded == value) return; _.unreaded = value; RaisePropertyChanged(); }
    }
    public int Page
    {
      get { return _page; }
      private set { Set(ref _page, value); }
    }
    public FriendlyDateTime PubDate
    {
      get
      {
        return new FriendlyDateTime(Utime.AddSeconds(_.update ));
      }
    }

    readonly string _category = null;
    private string first_or_default_category()
    {
      return _.categories.FirstOrDefault();
    }
    private bool append_category(string val)
    {
      if (val == first_or_default_category())
      {
        return false;
      }
      var n = new[] { val };
      _.categories = n.Concat(_.categories).ToArray();
      return true;
    }

    async void ExecuteDropSelf()
    {
      Debug.Assert(!string.IsNullOrEmpty(_.uri));
      //var rel = "/api/feed_source/unsubscribe.json?uri=" + System.Uri.EscapeDataString(_.uri);
      var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.FeedSourceUnsubscribe, new { _.uri });
      var s = await HttpClientUtils.Get<famousfront.datamodels.BackendError>(uri);
      var code = s.code != 0 ? s.code : s.data.code;
      var reason = s.code != 0 ? s.reason : s.data.reason;
      MessengerInstance.Send(new DropFeedSource() { model = this, code = code, reason = reason });
    }
    async void LoadUnreadCount()
    {      
      IsBusying = true;
      //var rel = "/api/feed_entry/source/unread_count.json?uri=" + System.Uri.EscapeDataString(_.uri);
      var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.FeedEntrySourceUnreadcount, new { _.uri });
      var s = await HttpClientUtils.Get<int>(uri);
      IsBusying = false;
      if (s.code != 0)
      {
        MessengerInstance.Send(new famousfront.messages.BackendError() { code = s.code, reason = s.reason });
        return;
      }
      UnreadCount = s.data;
    }

    internal void AddUnreadCount(int p)
    {
      UnreadCount += p;
    }

    void ExecuteGotoPage(int incre)
    {
      Page += incre;
      MessengerInstance.Send(this);
    }

  }
}
