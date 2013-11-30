using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using famousfront.datamodels;
using System.Windows.Input;
using GalaSoft.MvvmLight.Command;
using famousfront.utils;
using famousfront.messages;
using System.Diagnostics;
namespace famousfront.viewmodels
{
  internal class FeedSourceViewModel : famousfront.core.TaskViewModel
  {
    int _page;
    static readonly System.DateTime utime = new DateTime(1970, 1, 1, 0, 0, 0, 0);
    FeedSource _ = null;
    internal FeedSourceViewModel(FeedSource val)
    {
      _ = val;
      News = _.description;
      LoadUnreadCount();
    }
    public string Category { get { return _category ?? first_or_default_category(); } }
    public string Name { get { return _.name; } set { _.name = value; RaisePropertyChanged(); } }
    public string Uri { get { return _.uri; } }
    public string Description { get { return _.description; } }

    ICommand _goto_page;
    public ICommand GotoPageCommand
    {
      get { return _goto_page ?? (_goto_page = goto_page()); }
    }
    ICommand goto_page()
    {
      return new RelayCommand<int>(ExecuteGotoPage);
    }
    ICommand _drop_self;
    public ICommand DropSelfCommand
    {
      get { return _drop_self ?? (_drop_self = drop_self()); }
    }
    string _logo = null;
    public string Logo { get { return _logo; } private set { Set(ref _logo, value); } }

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
    public DateTime PubDate
    {
      get
      {
        return utime.AddSeconds(_.update );
      }
    }

    string _category = null;
    private string first_or_default_category()
    {
      return _.categories == null || _.categories.Length <= 0 ? "" : _.categories[0];
    }
    private bool append_category(string val)
    {
      if (val == first_or_default_category())
      {
        return false;
      }
      var n = new string[] { val };
      _.categories = n.Concat(_.categories).ToArray();
      return true;
    }
    private bool has_logo()
    {
      return false;
    }
    ICommand drop_self()
    {
      return new RelayCommand(ExecuteDropSelf);
    }
    async void ExecuteDropSelf()
    {
      Debug.Assert(!string.IsNullOrEmpty(_.uri));
      var rel = "/api/feed_source/unsubscribe.json?uri=" + System.Uri.EscapeDataString(_.uri);
      var s = await HttpClientUtils.Get<famousfront.datamodels.BackendError>(ServiceLocator.BackendPath(rel));
      var code = s.code != 0 ? s.code : s.data.code;
      var reason = s.code != 0 ? s.reason : s.data.reason;
      MessengerInstance.Send(new DropFeedSource() { model = this, code = code, reason = reason });
    }
    async void LoadUnreadCount()
    {      
      IsBusying = true;
      var rel = "/api/feed_entry/source/unread_count.json?uri=" + System.Uri.EscapeDataString(_.uri);
      var s = await HttpClientUtils.Get<int>(ServiceLocator.BackendPath(rel));
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
