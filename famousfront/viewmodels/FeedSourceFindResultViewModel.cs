using System;
using System.Collections.ObjectModel;
using System.ComponentModel;
using System.Windows.Data;
using famousfront.core;
using famousfront.datamodels;
using famousfront.utils;
using GalaSoft.MvvmLight.Threading;

namespace famousfront.viewmodels
{
  class FeedSourceFindResultViewModel : ViewModelBase
  {
    readonly ObservableCollection<FeedSourceFindEntryViewModel> _sources = new ObservableCollection<FeedSourceFindEntryViewModel>();
    ICollectionView _grouped_sources;
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
      //var rel = "/api/feed_source/find.json?q=" + Uri.EscapeDataString(q); //api/feed_source/find.json
      var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.FeedSourceFind, new {q});
      var v = await HttpClientUtils.Get<FeedEntity[]>(uri);
      if (v.code != 0)
      {
        MessengerInstance.Send(new BackendError { code = v.code, reason = v.reason });
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
}