using famousfront.core;
using famousfront.datamodels;
using System;
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
    public double Scale
    {
      get { return _.height > 0 ? (double)_.width / _.height : 0.0; }
    }
    public FeedMedia Self
    {
      get { return _; }
    }
    async Task DescribeImage()
    {
      if (_.width * _.height != 0 || string.IsNullOrEmpty(_.uri))
        return;
      //var rel = "/api/image/dimension.json?uri=" + Uri.EscapeDataString(_.uri);
      var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.ImageDimension, new { _.uri });

      var v = await famousfront.utils.HttpClientUtils.Get<FeedImage>(uri);
      if (v.code != 0)
      {
        _.duration = v.code;
        Reason = v.reason;
        MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
        return;
      }
      Width = v.data.width;
      Height = v.data.height;
      _.mime = v.data.mime;
      if (!string.IsNullOrEmpty(v.data.origin))
        _.local = v.data.origin;
      if (!string.IsNullOrEmpty(v.data.thumbnail))
        _.thumbanil = v.data.thumbnail;
      if (Width * Height != 0)
// ReSharper disable once ExplicitCallerInfoArgument
        RaisePropertyChanged("Scale");
    }
    async void LoadImage()
    {
      if (string.IsNullOrEmpty(_.uri))
        return;
      IsBusying = true;
      //var rel = "/api/image/description.json?uri=" + Uri.EscapeDataString(_.uri);
      var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.ImageDescription, new { _.uri });

      var v = await famousfront.utils.HttpClientUtils.Get<FeedImage>(uri);
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
    }
  }
}
