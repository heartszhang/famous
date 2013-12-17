using System;
using famousfront.datamodels;
using famousfront.utils;
using GalaSoft.MvvmLight.Threading;
using System.ComponentModel;
using System.Diagnostics;
using System.Windows.Data;
namespace famousfront.viewmodels
{
using FeedEntries = System.Collections.ObjectModel.ObservableCollection<FeedEntryViewModel>;
  class FeedEntriesViewModel : core.TaskViewModel
  {
    readonly int _page;
    readonly FeedSourceViewModel _parent;
    internal FeedEntriesViewModel(FeedSourceViewModel p)
    {
      _page = p.Page;
      _parent = p;
      Reload();
    }
    readonly FeedEntries _entries = new FeedEntries();
    ICollectionView _grouped_entries = null;
    public ICollectionView Entries { get { return _grouped_entries ?? (_grouped_entries = grouped_entries()); } }

    readonly VideoElementViewModel _video_service = new VideoElementViewModel();
    public core.TaskViewModel VideoService
    {
      get { return _video_service; }
    }
    ICollectionView grouped_entries()
    {
      var v = CollectionViewSource.GetDefaultView(_entries);
      v.GroupDescriptions.Add(new PropertyGroupDescription("PubDay"));
      return v;
    }
    async void Reload()
    {
      IsBusying = true;
      Debug.Assert(!string.IsNullOrEmpty(_parent.Uri));
      //var rel = "/api/feed_entry/unread.json?" + new { uri = _parent.Uri, count = 10, page = _page }.QueryString();
      var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.FeedEntryUnread, new { uri = _parent.Uri, count = 10, page = _page });
      var v = await HttpClientUtils.Get<FeedEntry[]>(uri);
      IsBusying = false;
      if (v.code != 0)
      {
        Reason = v.reason;
        MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
        return;
      }
      IsReady = true;
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _entries.Clear()), System.Windows.Threading.DispatcherPriority.ContextIdle);
      if (v.data == null)
        return;
      foreach (var fe in v.data)
      {
        var c = fe;
        await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _entries.Add(new FeedEntryViewModel(c))), System.Windows.Threading.DispatcherPriority.ContextIdle);
      }
    }
  }
}
