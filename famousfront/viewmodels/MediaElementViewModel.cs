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
    class MediaElementViewModel : ImageBaseViewModel
    {
        new FeedMedia _;
        internal MediaElementViewModel(FeedMedia v, FeedMedia backgroundimg) : base(backgroundimg)
        {
            _ = v;
        }
        public string VideoUrl
        {
            get { return _.uri; }
        }
        public new string Description
        {
            get { return _.description; }
        }
        ICommand _videoplay_command = null;
        public ICommand VideoPlayCommand
        {
            get { return _videoplay_command ?? (_videoplay_command = new RelayCommand(ExecuteVideoPlay)); }
        }
        void ExecuteVideoPlay()
        {
          MessengerInstance.Send(new VideoPlayRequest() { source = VideoUrl });
        }
    }
}
