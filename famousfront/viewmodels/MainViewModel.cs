using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using famousfront.core;
using famousfront.messages;
using GalaSoft.MvvmLight.Command;

namespace famousfront.viewmodels
{
    internal enum MainViewStatus
    {
        OFFLINE,
        READY,
        BOOTSTRAP,
        CONNECTING,
    }
    internal class MainViewModel : TaskViewModel
    {
        internal MainViewModel()
        {
            MessengerInstance.Register<BackendInitialized>(this, OnBackendInitialized);
            _content = new ContentViewModel();
            _offline = new OfflineViewModel();
            _booting = new BootstrapViewModel();
        }
        public override void Cleanup()
        {
            _content.Cleanup();
            _offline.Cleanup();
            _booting.Cleanup();
            base.Cleanup();
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

        BootstrapViewModel _booting;
        public BootstrapViewModel BootstrapViewModel
        {
            get { return _booting; }
            internal set
            {
                Set(ref _booting, value);
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

        MainViewStatus _status = MainViewStatus.BOOTSTRAP;
        public MainViewStatus Status
        {
            get { return _status; }
            internal set
            {
                Set(ref _status, value);
            }
        }

        void OnBackendInitialized(BackendInitialized msg)
        {
            Status = MainViewStatus.READY;
            _content.ReloadFeedSources();
        }
    }
}
