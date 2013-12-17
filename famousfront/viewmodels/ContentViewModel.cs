using System.Threading.Tasks;
using famousfront.datamodels;
using famousfront.messages;
using famousfront.utils;
using GalaSoft.MvvmLight.Command;
using System;
using System.Windows.Input;

namespace famousfront.viewmodels
{
  class ContentViewModel : core.TaskViewModel
  {
    readonly FeedSourcesViewModel _sources = new FeedSourcesViewModel();
    FeedEntriesViewModel _entries;
    readonly ImageTipViewModel _image_tip = new ImageTipViewModel();
    readonly System.Windows.Threading.DispatcherTimer _update_worker = new System.Windows.Threading.DispatcherTimer();
    bool _show_sources = true;
    public bool ShowFeedSources
    {
      get { return _show_sources; }
      internal set { Set(ref _show_sources, value); }
    }
    internal ContentViewModel()
    {
      _previous_entry_command = new RelayCommand(ExecutePreviousEntryCommand);
      _previous_source_command = new RelayCommand(ExecutePreviousSourceCommand);
      _next_entry_command = new RelayCommand(ExecuteNextEntryCommand);
      _next_source_command = new RelayCommand(ExecuteNextSourceCommand);

      MessengerInstance.Register<FeedSourceViewModel>(this, OnSelectedFeedSourceChanged);
      MessengerInstance.Register<ToggleFeedSource>(this, OnToggleFeedSource);
      _update_worker.Interval = TimeSpan.FromMinutes(ServiceLocator.FrontFlags.FeedUpdateInterval);
      _update_worker.Tick += do_updating;
      _update_worker.Start();
    }

    async void do_updating(object sender, EventArgs e)
    {
      var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.Tick);
      //const string rel = "/api/tick.json";
      var v = await HttpClientUtils.Get<BackendTick>(uri);
      if (v.code != 0)
      {
        MessengerInstance.Send(new famousfront.messages.BackendError() { code = v.code, reason = v.reason });
        return;
      }
      var bt = v.data.feeds;
      if (bt == null)
        return;
      foreach (var entity in bt)
      {
        MessengerInstance.Send(entity);
      }
    }
    public ImageTipViewModel ImageTipViewModel
    {
      get { return _image_tip; }
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
    RelayCommand<MouseButtonEventArgs> _toggle_feedsources_view_command;
    public ICommand ToggleFeedSourcesViewCommand
    {
      get { return _toggle_feedsources_view_command ?? (_toggle_feedsources_view_command = new RelayCommand<MouseButtonEventArgs>(ExecuteToggleFeedSources)); }
    }
    void ExecuteToggleFeedSources(MouseButtonEventArgs args)
    {
      OnToggleFeedSource(null);
    }
    void OnToggleFeedSource(ToggleFeedSource msg)
    {
      ShowFeedSources = !ShowFeedSources;
    }
    public override void Cleanup()
    {
      _update_worker.Stop();
      _update_worker.Tick -= do_updating;
      _toggle_feedsources_view_command = null;
      _previous_source_command = null;
      _previous_entry_command = null;
      _next_source_command = null;
      _next_entry_command = null;
      base.Cleanup();
    }
    internal async void ReloadFeedSources()
    {
      await _sources.Reload();
      LoadUpdates();
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
    async void LoadUpdates()
    {
      for (; ; )
      {
        IsBusying = true;
        //const string rel = "/api/update/popup.json";
        var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.UpdatePopup);

        var v = await HttpClientUtils.Get<famousfront.datamodels.FeedEntity>(uri);
        IsBusying = false;
        if (v.code != 0)
        {
          MessengerInstance.Send(new famousfront.messages.BackendError() { code = v.code, reason = v.reason });
          break;
        }
        var entity = v.data;
        if (entity == null)
          break;
        MessengerInstance.Send(entity);
        ServiceLocator.Log("pub-sub: {0}", entity.name??entity.uri);
        await Task.Delay(ServiceLocator.FrontFlags.KaPeriod * 7);
      }
    }
  }
}
