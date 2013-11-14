using famousfront.core;
using famousfront.datamodels;
using famousfront.messages;
using famousfront.utils;
using GalaSoft.MvvmLight.Command;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Input;

namespace famousfront.viewmodels
{
    class MediaElementViewModel : TaskViewModel
    {
        FeedMedia _;
        FeedMedia _background;
        internal MediaElementViewModel(FeedMedia v, FeedMedia backgroundimg)
        {
            _ = v;
            _background = backgroundimg;
            if (_background != null)
                LoadImage();
        }
        public string Url
        {
            get { return _.uri; }
        }
        public string Description
        {
            get { return _.description; }
        }
        string _background_image;
        public string BackgroundImage
        {
            get { return _background_image; }
            private set { Set(ref _background_image, value); }
        }
        ICommand _videoplay_command = null;
        public ICommand VideoPlayCommand
        {
            get { return _videoplay_command ?? (_videoplay_command = new RelayCommand<string>(ExecuteVideoPlay)); }
        }
        void ExecuteVideoPlay(string url)
        {
            MessengerInstance.Send(new VideoPlayRequest() { source = url});
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
                MessengerInstance.Send(new famousfront.messages.BackendError() { code = v.code, reason = v.reason });
                return;
            }
            IsReady = true;
            _.width = v.data.width;
            _.height = v.data.height;
            _.mime = v.data.mime;
            _.local = v.data.origin;
            _.thumbanil = v.data.thumbnail;
            BackgroundImage = _.thumbanil;
        }
    }
}
