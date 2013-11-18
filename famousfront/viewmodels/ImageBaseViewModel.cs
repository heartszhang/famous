using famousfront.core;
using famousfront.datamodels;
using GalaSoft.MvvmLight.Threading;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.viewmodels
{
  abstract class ImageBaseViewModel : TaskViewModel
  {
    protected FeedMedia _;
    internal ImageBaseViewModel(FeedMedia v)
    {
      _ = v;
      LoadImage();
    }
    string _ideal_url;
    public string IdealUrl
    {
      get { return _ideal_url; }
      protected set { Set(ref _ideal_url, value); }
    }
    public string Url
    {
      get { return _.thumbanil; }
      protected set { _.thumbanil = value; RaisePropertyChanged(); }
    }
    public string OriginUrl
    {
      get { return _.local; }
      protected set { _.local = value; RaisePropertyChanged(); }
    }
    public string Description
    {
      get { return _.description; }
      protected set { _.description = value; RaisePropertyChanged(); }
    }
    public int Width
    {
      get { return _.width; }
      protected set { _.width = value; RaisePropertyChanged(); }
    }
    public int Height
    {
      get { return _.height; }
      protected set { _.height = value; RaisePropertyChanged(); }
    }
    async void LoadImage()
    {
      IsBusying = true;
      var rel = "/api/image/description.json?uri=" + Uri.EscapeDataString(_.uri);
      var v = await famousfront.utils.HttpClientUtils.Get<FeedImage>(ServiceLocator.BackendPath(rel));
      IsBusying = false;
      if (v.code != 0)
      {
        Reason = v.reason;
        MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
        return;
      }
      IsReady = true;
      Width = v.data.width;
      Height = v.data.height;
      _.mime = v.data.mime;
      OriginUrl = v.data.origin;
      Url = v.data.thumbnail;
      var scale = Height > 0 ? (Width * 100 / Height) : 0;
      IdealUrl = (scale >= 100) ? OriginUrl : Url;
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
      }), System.Windows.Threading.DispatcherPriority.ContextIdle); 
    }
  }
}
