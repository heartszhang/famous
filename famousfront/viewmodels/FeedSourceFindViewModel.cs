using famousfront.core;
using GalaSoft.MvvmLight.Command;
using System;
using System.Windows.Input;

namespace famousfront.viewmodels
{
  class FeedSourceFindViewModel : TaskViewModel
  {
    ViewModelBase _content;
    string _query;
    readonly ICommand _feedsource_find_command;

    internal FeedSourceFindViewModel()
    {
      _feedsource_find_command = new RelayCommand<string>(ExecuteFindFeedSource);
    }
    public string Query { get { return _query; } set { Set(ref _query, value); } }

    public ViewModelBase Content
    {
      get { return _content; }
      set { Set(ref _content, value); }
    }

    public ICommand FeedSourceFindCommand
    {
      get { return _feedsource_find_command ; }
    }

    async void ExecuteShowFeedSourceImp(string q)
    {
      IsBusying = true;
      var cnt = new FeedSourceShowResultViewModel();
      Reason = await cnt.Load(q);
      Content = cnt;
      IsBusying = false;
      IsReady = string.IsNullOrEmpty(Reason);
    }
    async void ExecuteFindFeedSourceImp(string q)
    {
      IsBusying = true;
      var cnt = new FeedSourceFindResultViewModel();
      Reason = await cnt.Load(q);
      IsReady = string.IsNullOrEmpty(Reason);
      Content = cnt;
      IsBusying = false;
    }
    void ExecuteFindFeedSource(string q)
    {
      Uri u;
      if (Uri.TryCreate(q, UriKind.Absolute, out u))
        ExecuteShowFeedSourceImp(q);
      else ExecuteFindFeedSourceImp(q);
    }
  }
}
