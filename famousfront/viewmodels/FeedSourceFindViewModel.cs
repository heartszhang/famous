using famousfront.core;
using famousfront.datamodels;
using famousfront.utils;
using GalaSoft.MvvmLight.Command;
using System;
using System.Windows.Input;
  using System.ComponentModel;
  using System.Windows.Data;
using GalaSoft.MvvmLight.Threading;

namespace famousfront.viewmodels
{
  using FeedEntries = System.Collections.ObjectModel.ObservableCollection<FeedEntryViewModel>;
  using FeedSources = System.Collections.ObjectModel.ObservableCollection<FeedSourceFindEntryViewModel>;
  using System.Diagnostics;
  class FeedSourceFindResultViewModel : ViewModelBase
  {
    readonly FeedSources _sources = new FeedSources();
    ICollectionView _grouped_sources = null;
    public ICollectionView Sources
    {
      get { return _grouped_sources ?? (_grouped_sources = grouped_sources()); }
    }

    ICollectionView grouped_sources()
    {
      var v = CollectionViewSource.GetDefaultView(_sources);
      return v;
    }
    internal async System.Threading.Tasks.Task<string> Load(string q)
    {
      var rel = "/api/feed_source/find.json?q=" + Uri.EscapeDataString(q); //api/feed_source/find.json
      var v = await HttpClientUtils.Get<FeedSourceFindEntry[]>(ServiceLocator.BackendPath(rel));
      if (v.code != 0)
      {
        MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
        return v.reason;
      }
      var fss = v.data;
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _sources.Clear()), System.Windows.Threading.DispatcherPriority.ContextIdle);
      foreach (var f in fss)
      {
        var c = f;
        await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _sources.Add(new FeedSourceFindEntryViewModel(c))), System.Windows.Threading.DispatcherPriority.ContextIdle);
      }
      return string.Empty;
    }
  }
  class FeedSourceShowResultViewModel : TaskViewModel
  {
    FeedEntity _;

    readonly FeedEntries _entries = new FeedEntries();
    ICollectionView _grouped_entries = null;
    public ICollectionView Entries { get { return _grouped_entries ?? (_grouped_entries = grouped_entries()); } }
    ICollectionView grouped_entries()
    {
      var v = CollectionViewSource.GetDefaultView(_entries);
      return v;
    }
    internal async System.Threading.Tasks.Task<string> Load(string q)
    {
      var rel = "/api/feed_source/show.json?q=" + Uri.EscapeDataString(q); //api/feed_source/find.json
      var v = await HttpClientUtils.Get<FeedEntity>(ServiceLocator.BackendPath(rel));
      if (v.code != 0)
      {
        MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
        return v.reason;
      }
      _ = v.data;
      RaisePropertyChanged("Name");
      foreach (var e in _.entries)
      {
        var c = e;
        await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _entries.Add(new FeedEntryViewModel(c))), System.Windows.Threading.DispatcherPriority.ContextIdle);
      }
      return string.Empty;
    }
    public string Name
    {
      get { return _ == null ? "" : _.name; }
    }
    bool _subscribed;
    public bool HasSubscribed { get { return _subscribed; } set { _subscribed = value; RaisePropertyChanged(); } }
    ICommand _subscribe_self;
    public ICommand SubscribeCommand
    {
      get { return _subscribe_self ?? (_subscribe_self = subscribe_self()); }
    }
    ICommand subscribe_self()
    {
      return new RelayCommand<bool>(ExecuteSubscribeSelf);
    }
    async void ExecuteSubscribeSelf(bool sub)
    {
      if ((sub && HasSubscribed) || (!sub && !HasSubscribed))
        return;
      IsBusying = true;
      if (sub)
      {
        Debug.Assert(!string.IsNullOrEmpty(_.uri));
        var rel = "/api/feed_source/subscribe.json?uri=" + Uri.EscapeDataString(_.uri);
        var v = await famousfront.utils.HttpClientUtils.Get<FeedSource>(ServiceLocator.BackendPath(rel));
        if (v.code != 0)
        {
          Reason = v.reason;
          MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
          return;
        }
        MessengerInstance.Send(new messages.SubscribeFeedSource { source = v.data });
        HasSubscribed = true;
      }
      else
      {
        MessengerInstance.Send(new messages.UnsubscribeFeedSource { source = _.uri });
      }
      IsBusying = false;
      /*	http.HandleFunc("/api/feed_source/unsubscribe.json", webapi_feedsource_unsubscribe)*/
    }
  }
  class FeedSourceFindViewModel : TaskViewModel
  {
    string _query;
    public string Query { get { return _query; } set { Set(ref _query, value); } }

    ViewModelBase _content;
    public ViewModelBase Content
    {
      get { return _content; }
      set { Set(ref _content, value); }
    }
    ICommand _feedsource_find_command;
    public ICommand FeedSourceFindCommand
    {
      get { return _feedsource_find_command ?? (_feedsource_find_command = feedsource_find_command()); }
    }
    ICommand feedsource_find_command()
    {
      return new RelayCommand<string>(ExecuteFindFeedSource);
    }
    async void ExecuteShowFeedSourceImp(string q)
    {
      IsBusying = true;
      var cnt = new FeedSourceShowResultViewModel();
      Reason = await cnt.Load(q);
      Content = cnt;
      IsBusying = false;
      IsReady = string.IsNullOrEmpty(Reason);
    }
    async void ExecuteFindFeedSourceImp(string q)
    {
      IsBusying = true;
      var cnt = new FeedSourceFindResultViewModel();
      Reason = await cnt.Load(q);
      IsReady = string.IsNullOrEmpty(Reason);
      Content = cnt;
      IsBusying = false;
    }
    void ExecuteFindFeedSource(string q)
    {
      Uri u;
      if (Uri.TryCreate(q, UriKind.Absolute, out u))
      {
        ExecuteShowFeedSourceImp(q);
        return;
      }
      else
      {
        ExecuteFindFeedSourceImp(q);
      }
    }
  }
}
