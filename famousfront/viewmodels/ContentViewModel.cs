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
    class ContentViewModel : famousfront.core.ViewModelBase
    {
        FeedSourcesViewModel _sources = new FeedSourcesViewModel();
        FeedEntriesViewModel _entries = null;
        bool _show_sources = true;
        public bool ShowFeedSources
        {
            get { return _show_sources; }
            internal set { Set(ref _show_sources, value); }
        }
        internal ContentViewModel()
        {
            _toggle_feedsources_command = new RelayCommand<MouseButtonEventArgs>(ExecuteToggleFeedSources);
            _previous_entry_command = new RelayCommand(ExecutePreviousEntryCommand);
            _previous_source_command = new RelayCommand(ExecutePreviousSourceCommand);
            _next_entry_command = new RelayCommand(ExecuteNextEntryCommand);
            _next_source_command = new RelayCommand(ExecuteNextSourceCommand);

            MessengerInstance.Register<FeedSourceViewModel>(this, OnSelectedFeedSourceChanged);
            MessengerInstance.Register<ToggleFeedSource>(this,OnToggleFeedSource);
        }
        public FeedEntriesViewModel FeedEntriesViewModel
        {
            get { return _entries; }
            internal set { Set(ref _entries, value); }
        }
        public FeedSourcesViewModel FeedSourcesViewModel
        {
            get { return _sources; }
        }
        RelayCommand<MouseButtonEventArgs> _toggle_feedsources_command;
        public ICommand ToggleFeedSourcesCommand
        {
            get { return _toggle_feedsources_command; }
        }
        void ExecuteToggleFeedSources(MouseButtonEventArgs args)
        {
            OnToggleFeedSource(null);
        }
        void OnToggleFeedSource(ToggleFeedSource msg){
            ShowFeedSources = !ShowFeedSources;            
        }
        public override void Cleanup()
        {
            _toggle_feedsources_command = null;
            _previous_source_command = null;
            _previous_entry_command = null;
            _next_source_command = null;
            _next_entry_command = null;
            base.Cleanup();
        }
        internal void ReloadFeedSources()
        {
            _sources.Reload();
        }

        RelayCommand _previous_entry_command;
        public RelayCommand PreviousEntryCommand
        {
            get { return _previous_entry_command; }
        }
        RelayCommand _next_entry_command;
        public RelayCommand NextEntryCommand
        {
            get { return _next_entry_command; }
        }

        RelayCommand _previous_source_command;
        public RelayCommand PreviousSourceCommand
        {
            get { return _previous_source_command; }
        }
        RelayCommand _next_source_command;
        public RelayCommand NextSourceCommand
        {
            get { return _next_source_command; }
        }
        void ExecutePreviousEntryCommand()
        {

        }
        void ExecutePreviousSourceCommand()
        {

        }
        void ExecuteNextEntryCommand()
        {

        }
        void ExecuteNextSourceCommand()
        {

        }
        void OnSelectedFeedSourceChanged(FeedSourceViewModel cur)
        {
            if (FeedEntriesViewModel != null)
            {
                FeedEntriesViewModel.Cleanup();
                FeedEntriesViewModel = null;
            }
            if (cur == null)
            {
                return;
            }
            FeedEntriesViewModel = new FeedEntriesViewModel(cur);
        }
    }
}
