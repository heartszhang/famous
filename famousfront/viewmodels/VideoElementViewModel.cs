using famousfront.core;
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
    class VideoElementViewModel : TaskViewModel
    {
        string _current_source;
        public string CurrentSource
        {
            get { return _current_source; }
            private set { Set(ref _current_source, value); }
        }
        internal VideoElementViewModel()
        {
            MessengerInstance.Register<VideoPlayRequest>(this, OnVideoPlayRequest);
        }
        void OnVideoPlayRequest(VideoPlayRequest msg)
        {
            CurrentSource = msg.source;
        }
        ICommand _videoplay_command = null;
        public ICommand VideoStopCommand
        {
            get { return _videoplay_command ?? (_videoplay_command = new RelayCommand(ExecuteVideoStop)); }
        }
        void ExecuteVideoStop()
        {
            MessengerInstance.Send(new VideoPlayRequest() { source = null });
        }
    }
}
