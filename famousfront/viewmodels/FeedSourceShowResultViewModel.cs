using System;
using System.Collections.ObjectModel;
using System.ComponentModel;
using System.Diagnostics;
using System.Windows.Data;
using System.Windows.Input;
using famousfront.core;
using famousfront.datamodels;
using famousfront.utils;
using GalaSoft.MvvmLight.Command;
using GalaSoft.MvvmLight.Threading;

namespace famousfront.viewmodels
{
  class FeedSourceShowResultViewModel : TaskViewModel
  {
    FeedEntity      _;
    bool            _subscribed;
    ICollectionView _grouped_entries;

    readonly ObservableCollection<FeedEntryViewModel> _entries = new ObservableCollection<FeedEntryViewModel>();
    public ICollectionView Entries { get { return _grouped_entries ?? (_grouped_entries = grouped_entries()); } }
    ICollectionView grouped_entries()
    {
      var v = CollectionViewSource.GetDefaultView(_entries);
      return v;
    }
    internal async System.Threading.Tasks.Task<string> Load(string q)
    {
      //var rel = "/api/feed_source/show.json?q=" + Uri.EscapeDataString(q); //api/feed_source/find.json
      var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.FeedSourceShow, new { q });
      var v = await HttpClientUtils.Get<FeedEntity>(uri);
      if (v.code != 0)
      {
        MessengerInstance.Send(new BackendError { code = v.code, reason = v.reason });
        return v.reason;
      }
      _ = v.data;      
      Name = v.data.name;
      HasSubscribed = v.data.subscribe_state == FeedSourceSubscribeStates.Subscribed;
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
      set { _.name = value; RaisePropertyChanged(); }
    }
    public bool HasSubscribed { get { return _subscribed; } set { Set(ref _subscribed , value);} }
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
        //var rel = "/api/feed_source/subscribe.json?uri=" + Uri.EscapeDataString(_.uri);
        var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.FeedSourceSubscribe, new { _.uri });
        var v = await HttpClientUtils.Get<FeedSource>(uri);
        if (v.code != 0)
        {
          Reason = v.reason;
          MessengerInstance.Send(new BackendError { code = v.code, reason = v.reason });
          return;
        }
        MessengerInstance.Send(new messages.SubscribeFeedSource { source = v.data });
      }
      else
      {
        MessengerInstance.Send(new messages.UnsubscribeFeedSource { source = _.uri });
      }
      HasSubscribed = sub;
      IsBusying = false;
    }
  }
}