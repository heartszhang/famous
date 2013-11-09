using famousfront.core;
using famousfront.datamodels;
using famousfront.utils;
using GalaSoft.MvvmLight.Threading;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.viewmodels
{
    class ImageElementViewModel : TaskViewModel
    {
        FeedMedia _;
        internal ImageElementViewModel(FeedMedia v)
        {
            _ = v;
            LoadImage();
        }
        string _url;
        public string Url
        {
            get { return _url; }
            internal set{ Set(ref _url, value);}
        }
        public string Description
        {
            get { return _.description; }
        }
        async void LoadImage()
        {
            IsBusying = true;
            var rel = "/api/image/description.json?uri=" + Uri.EscapeDataString(_.uri);
            var v = await HttpClientUtils.Get<FeedImage>(ServiceLocator.BackendPath(rel));
            IsBusying = false;
            if (v.code != 0)
            {
                Reason = v.reason;
                MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
                return;
            }
            IsReady = true;
            _.width = v.data.width;
            _.height = v.data.height;
            _.mime = v.data.mime;
            _.local = v.data.origin;
            _.thumbanil = v.data.thumbnail;
            Url = _.thumbanil;            
           // RaisePropertyChanged("Url");
//            await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
//            {
//            }), System.Windows.Threading.DispatcherPriority.ContextIdle); 
        }
    }
}
