using famousfront.core;
using famousfront.datamodels;
using famousfront.utils;
using GalaSoft.MvvmLight.Threading;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Data;

namespace famousfront.viewmodels
{
  using FeedImages = System.Collections.ObjectModel.ObservableCollection<ImageUnitViewModel>;

  class ImageGalleryViewModel : TaskViewModel
  {
    FeedMedia[] _;
    internal ImageGalleryViewModel(FeedMedia[] imgs)
    {
      _ = imgs;
      LoadImages();
    }
    FeedImages _coll_images = null;
    ICollectionView _images = null;
    public ICollectionView Images
    {
      get { return _images; }
      private set { Set(ref _images, value); }
    }
    async void LoadImages()
    {
      if (_ == null || _.Length == 0)
        return;
      IsBusying = true;

      await Task.WhenAll(_.Select(m => DescribeImage(m))).ConfigureAwait(false);
      IsBusying = false;

      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        _coll_images = new FeedImages();
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
      foreach (var i in _)
      {
        var c = i;
        if (c.duration != 0)
          continue;
        await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
        {
          _coll_images.Add(new ImageUnitViewModel(c));
        }), System.Windows.Threading.DispatcherPriority.ContextIdle);
      }
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        Images = CollectionViewSource.GetDefaultView(_coll_images);
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);

      IsReady = true;
    }
    async Task DescribeImage(FeedMedia img)
    {
      var rel = "/api/image/dimension.json?uri=" + Uri.EscapeDataString(img.uri);
//      var rel = "/api/image/description.json?uri=" + Uri.EscapeDataString(img.uri);
      var v = await HttpClientUtils.Get<FeedImage>(ServiceLocator.BackendPath(rel));
      if (v.code != 0)
      {
        img.duration = v.code;
        Reason = v.reason;
        MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
        return;
      }
      img.width = v.data.width;
      img.height = v.data.height;
      img.mime = v.data.mime;
      img.local = v.data.origin;
      img.thumbanil = v.data.thumbnail;

      // RaisePropertyChanged("Url");
      //            await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      //            {
      //            }), System.Windows.Threading.DispatcherPriority.ContextIdle); 
    }

  }
}
