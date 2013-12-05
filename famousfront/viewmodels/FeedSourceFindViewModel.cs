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
  using FeedSources = System.Collections.ObjectModel.ObservableCollection<FeedSourceFindEntryViewModel>;
  class FeedSourceFindViewModel : TaskViewModel
  {
    string _query;
    public string Query { get { return _query; } set { Set(ref _query, value); } }
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
    ICommand _feedsource_find_command;
    public ICommand FeedSourceFindCommand
    {
      get { return _feedsource_find_command ?? (_feedsource_find_command = feedsource_find_command()); }
    }
    ICommand feedsource_find_command()
    {
      return new RelayCommand<string>(ExecuteFindFeedSource);
    }
    async void ExecuteFindFeedSource(string q)
    {
      IsBusying = true;
      var rel = "/api/feed_source/find.json?q=" + Uri.EscapeDataString(q); //api/feed_source/find.json
      var v = await HttpClientUtils.Get<FeedSourceFindEntry[]>(ServiceLocator.BackendPath(rel));
      IsBusying = false;
      if (v.code != 0)
      {
        Reason = v.reason;
        MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
        return;
      }
      var fss = v.data;
      IsReady = true;
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _sources.Clear()), System.Windows.Threading.DispatcherPriority.ContextIdle);
      foreach (var f in fss)
      {
        var c = f;
        await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _sources.Add(new FeedSourceFindEntryViewModel(c))), System.Windows.Threading.DispatcherPriority.ContextIdle);
      }
    }
  }
}
