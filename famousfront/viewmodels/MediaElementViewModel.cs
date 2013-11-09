using famousfront.core;
using famousfront.datamodels;
using famousfront.messages;
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
        internal MediaElementViewModel(FeedMedia v)
        {
            _ = v;
        }
        public string Url
        {
            get { return _.uri; }
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
    }
}
