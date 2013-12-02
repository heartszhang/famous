using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using famousfront.datamodels;
using famousfront.utils;
using GalaSoft.MvvmLight.Threading;
using System.ComponentModel;
using System.Diagnostics;
using System.Net;
using System.Windows.Data;
namespace famousfront.viewmodels
{
using FeedEntries = System.Collections.ObjectModel.ObservableCollection<FeedEntryViewModel>;
  class FeedEntriesViewModel : famousfront.core.TaskViewModel
  {
    int _page;
    FeedSourceViewModel _parent;
    FeedEntries _entries = new FeedEntries();
    internal FeedEntriesViewModel(FeedSourceViewModel p)
    {
      _page = p.Page;
      _parent = p;
      Reload();
    }
    ICollectionView _grouped_entries = null;
    public ICollectionView Entries { get { return _grouped_entries ?? (_grouped_entries = grouped_entries()); } }

    VideoElementViewModel _video_service = new VideoElementViewModel();
    public famousfront.core.TaskViewModel VideoService
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
      var rel = "/api/feed_entry/unread.json?" + new { uri = _parent.Uri, count = 10, page = _page }.QueryString();
      //var rel = "/api/feed_entry/unread.json?uri=" + Uri.EscapeDataString(_parent.Uri);
      var v = await HttpClientUtils.Get<FeedEntry[]>(ServiceLocator.BackendPath(rel));
      IsBusying = false;
      if (v.code != 0)
      {
        Reason = v.reason;
        MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
        return;
      }
      IsReady = true;
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        _entries.Clear();
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
      if (v.data == null)
        return;
      foreach (var fe in v.data)
      {
        var c = fe;
        await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
        {
          _entries.Add(new FeedEntryViewModel(c));
        }), System.Windows.Threading.DispatcherPriority.ContextIdle);
      }
    }
  }
}
