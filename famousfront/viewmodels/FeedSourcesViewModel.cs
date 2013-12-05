using System;
using System.Collections.ObjectModel;
using System.Linq;
using System.Threading.Tasks;
using System.ComponentModel;
using System.Windows.Data;
using GalaSoft.MvvmLight.Threading;
using famousfront.utils;
using famousfront.messages;
using famousfront.datamodels;

namespace famousfront.viewmodels
{
  using FeedSources = ObservableCollection<FeedSourceViewModel>;

  class FeedSourcesViewModel : core.TaskViewModel
  {
    internal FeedSourcesViewModel()
    {
      MessengerInstance.Register<DropFeedSource>(this, OnDropFeedSource);
      MessengerInstance.Register<SubscribeFeedSource>(this, OnSubscribeFeedSource);
      MessengerInstance.Register<UnsubscribeFeedSource>(this, OnUnsubscribeFeedSource);
      MessengerInstance.Register<FeedEntity>(this, OnFeedEntity);
    }

    private void OnFeedEntity(FeedEntity obj)
    {
      var s = _sources.FirstOrDefault(v => obj.uri == v.Uri);
      if (s == null)
        return;
      s.AddUnreadCount(obj.entries.Length);
    }

    private void OnDropFeedSource(DropFeedSource obj)
    {
      if (obj.code != 0)
        return;
      DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _sources.Remove(obj.model)));
    }
    private async void OnSubscribeFeedSource(SubscribeFeedSource msg)
    {
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _sources.Add(new FeedSourceViewModel(msg.source))), System.Windows.Threading.DispatcherPriority.ContextIdle);
    }
    private async void OnUnsubscribeFeedSource(UnsubscribeFeedSource msg)
    {
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        var vm = _sources.First(s => s.Uri == msg.source);
        if (vm != null)
          _sources.Remove(vm);
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);

    }

    readonly FeedSources _sources = new FeedSources();
    ICollectionView _grouped_sources = null;
    public ICollectionView Sources
    {
      get { return _grouped_sources ?? (_grouped_sources = grouped_sources()); }
    }
    int _selected_index = -1;
    public int SelectedIndex
    {
      get { return _selected_index; }
      set { Set(ref _selected_index, value); }
    }
    FeedSourceViewModel _selected;
    public FeedSourceViewModel Selected
    {
      get { return _selected; }
      set
      {
        var prev = _selected;
        Set(ref _selected, value);
        if (prev != value)
          MessengerInstance.Send(value);
      }
    }

    ICollectionView grouped_sources()
    {
      var v = CollectionViewSource.GetDefaultView(_sources);
      v.GroupDescriptions.Add(new PropertyGroupDescription("Category"));
      return v;
    }
    internal async Task Reload()
    {
      IsBusying = true;
      var fs = await HttpClientUtils.Get<FeedSource[]>(ServiceLocator.BackendPath("/api/feed_source/all.json"));
      IsBusying = false;
      if (fs.code != 0)
      {
        MessengerInstance.Send(new messages.BackendError() { code = fs.code, reason = fs.reason });
        return;
      }
      var fss = fs.data;
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _sources.Clear()), System.Windows.Threading.DispatcherPriority.ContextIdle);
      foreach (var f in fss)
      {
        var c = f;
        await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _sources.Add(new FeedSourceViewModel(c))), System.Windows.Threading.DispatcherPriority.ContextIdle);
      }
    }
  }
}
