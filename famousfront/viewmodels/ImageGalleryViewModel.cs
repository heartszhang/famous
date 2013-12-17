using famousfront.core;
using famousfront.datamodels;
using famousfront.utils;
using GalaSoft.MvvmLight.Threading;
using System;
using System.ComponentModel;
using System.Linq;
using System.Threading.Tasks;
using System.Windows.Data;
using GalaSoft.MvvmLight.Command;
using System.Windows.Input;

namespace famousfront.viewmodels
{
  using FeedImages = System.Collections.ObjectModel.ObservableCollection<ImageUnitViewModel>;
  class ImagePanelViewModel : TaskViewModel
  {
    bool _show_panel;
    bool _initialized;
    readonly ICommand _toggle_show_panel;
    readonly FeedMedia[] _;

    public ImagePanelViewModel(FeedMedia[] imgs)
    {
      _ = imgs;
      _toggle_show_panel = new RelayCommand(ExecuteToggleShowPanel);
      DescribeImages();
    }
    FeedImages _coll_images = null;
    ICollectionView _images = null;
    public ICollectionView Images
    {
      get { return _images; }
      private set { Set(ref _images, value); }
    }
    async void DescribeImages()
    {
      IsBusying = true;
      await Task.WhenAll(_.Select(DescribeImage)).ConfigureAwait(false);
      await LoadImages();
      IsBusying = false;
    }
    internal async Task LoadImages()
    {
      if (_ == null || _.Length == 0)
        return;
      if (_initialized)
        return;
      _initialized = true;
      IsBusying = true;

      await Task.WhenAll(_.Select(DescribeImage)).ConfigureAwait(false);
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
        await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() => _coll_images.Add(new ImageUnitViewModel(c))), System.Windows.Threading.DispatcherPriority.ContextIdle);
      }
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        Images = CollectionViewSource.GetDefaultView(_coll_images);
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);

      IsReady = true;
    }
    async Task DescribeImage(FeedMedia img)
    {
      if (img.width * img.height != 0)
        return;
      //var rel = "/api/image/dimension.json?uri=" + Uri.EscapeDataString(img.uri);
      var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.ImageDimension, new { img.uri });
      var v = await HttpClientUtils.Get<FeedImage>(uri);
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
    }
    public bool IsShowPanel
    {
      get { return _show_panel; }
      protected set { Set(ref _show_panel, value); }
    }

    public ICommand ToggleShowPanelCommand
    {
      get { return _toggle_show_panel; }
    }

    void ExecuteToggleShowPanel()
    {
      IsShowPanel = !IsShowPanel;
    }
  }
}
