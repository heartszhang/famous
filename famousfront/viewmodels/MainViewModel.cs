using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using famousfront.core;

namespace famousfront.viewmodels
{
    internal enum MainViewStatus
    {
        OFFLINE,
        READY,
        CONNECTING,
    }
    internal class MainViewModel : TaskViewModel
    {
        internal MainViewModel()
        {
            _content = new ContentViewModel();
            _offline = new OfflineViewModel();
        }
        ContentViewModel _content;
        public ContentViewModel ContentViewModel
        {
            get { return _content; }
            internal set
            {
                Set(ref _content, value);
            }
        }

        ConnectingViewModel _connecting;
        public ConnectingViewModel ConnectingViewModel
        {
            get { return _connecting; }
            internal set
            {
                Set(ref _connecting, value);
            }
        }

        OfflineViewModel _offline;
        public OfflineViewModel OfflineViewModel
        {
            get { return _offline; }
            internal set
            {
                Set(ref _offline, value);
            }
        }
        SettingsViewModel _settings;
        public SettingsViewModel SettingsViewModel
        {
            get { return _settings; }
            internal set
            {
                Set(ref _settings, value);
            }
        }

        bool _show_settings = false;
        public bool ShowSettings
        {
            get { return _show_settings; }
            internal set 
            {
                Set(ref _show_settings, value);
            }
        }

        MainViewStatus _status = MainViewStatus.READY;
        public MainViewStatus Status
        {
            get { return _status; }
            internal set
            {
                Set(ref _status, value);
            }
        }
    }
}
