using System;
using System.Diagnostics;
using System.Windows.Input;
using famousfront.messages;
using GalaSoft.MvvmLight.Command;
using GalaSoft.MvvmLight.Messaging;

namespace famousfront
{
  class CommandsService
  {
    readonly ICommand _hyperlink_navigate;
    readonly ICommand _show_messages_view;
    readonly ICommand _show_find_feedsource_view;
    readonly ICommand _toggle_feedsources_view;
    protected CommandsService()
    {
      _toggle_feedsources_view = new RelayCommand(DispatchToggleFeedSource);
      _show_find_feedsource_view = new RelayCommand(DispatchShowFindFeedSourceView);
      _show_messages_view = new RelayCommand(DispatchShowMessageView);
      _hyperlink_navigate = new RelayCommand<Uri>(DispatchHyperlinkNavigate);
    }

    public ICommand ToggleFeedSourcesViewCommand
    {
      get { return _toggle_feedsources_view ; }
    }

    void DispatchToggleFeedSource()
    {
      Messenger.Default.Send(new ToggleFeedSource());
    }

    public ICommand FindFeedSourceCommand
    {
      get { return _show_find_feedsource_view ; }
    }

    void DispatchShowFindFeedSourceView()
    {
      Messenger.Default.Send(new ShowFindFeedSourceView());
    }

    public ICommand ShowMessagesCommand
    {
      get { return _show_messages_view ; }
    }

    void DispatchShowMessageView()
    {
      Messenger.Default.Send(new ShowMessagesView());
    }

    public ICommand HyperlinkNavigateCommand
    {
      get { return _hyperlink_navigate ; }
    }

    void DispatchHyperlinkNavigate(Uri url)
    {
      if (url == null)
        return;
      using (var p = Process.Start(url.ToString())) { };
    }
  }
}