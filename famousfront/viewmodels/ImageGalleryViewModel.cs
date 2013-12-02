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
  using GalaSoft.MvvmLight.Command;
  using System.Windows.Input;
  using FeedImages = System.Collections.ObjectModel.ObservableCollection<ImageUnitViewModel>;
  class ImagePanelViewModel : TaskViewModel
  {
    FeedMedia[] _;

    public ImagePanelViewModel(FeedMedia[] imgs)
    {
      _ = imgs;
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
      await Task.WhenAll(_.Select(m => DescribeImage(m))).ConfigureAwait(false);
      IsBusying = false;
    }
    bool _initialized;
    internal async void LoadImages()
    {
      if (_ == null || _.Length == 0)
        return;
      if (_initialized)
        return;
      _initialized = true;
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
      if (img.width * img.height != 0)
        return;
      var rel = "/api/image/dimension.json?uri=" + Uri.EscapeDataString(img.uri);
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
    }
  }
  class ImageGalleryViewModel : TaskViewModel
  {
    internal ImageGalleryViewModel(FeedMedia[] imgs) 
    {
      _panel = new ImagePanelViewModel(imgs);
      if (imgs[0] != null)
        _first = new ImageElementViewModel(imgs[0]);
      if (imgs.Length < ServiceLocator.Flags.ShowGalleryThreshold)
        ExecuteToggleShowPanel();
    }
    ImagePanelViewModel _panel;
    public ImagePanelViewModel ImagePanelViewModel
    {
      get { return _panel; }
    }
    ImageElementViewModel _first;
    public ImageElementViewModel ImageElementViewModel
    {
      get { return _first; }
    }
    bool _show_panel;
    public bool IsShowPanel 
    {
      get { return _show_panel; }
      protected set { Set(ref _show_panel, value); }
    }
    ICommand _toggle_show_panel;
    public ICommand ToggleShowPanelCommand
    {
      get { return _toggle_show_panel ?? (_toggle_show_panel = toggle_show_panel());}
    }
    ICommand toggle_show_panel()
    {
      return new RelayCommand(ExecuteToggleShowPanel);
    }
    bool _panel_initialized;
    void ExecuteToggleShowPanel()
    {
      IsShowPanel = !IsShowPanel;
      if (!_panel_initialized)
      {
        _panel_initialized = true;
        ImagePanelViewModel.LoadImages();
      }
    }
  }
}
