using famousfront.core;
using famousfront.datamodels;
using GalaSoft.MvvmLight.Command;
using System.Diagnostics;
using System.Windows.Input;

namespace famousfront.viewmodels
{
  class FeedSourceFindEntryViewModel : TaskViewModel
  {
    bool _has_subscribed;
    readonly ICommand _subscribe_self;
    readonly FeedEntity _;
    internal FeedSourceFindEntryViewModel(FeedEntity v)
    {
      _ = v;
      _has_subscribed = _.subscribe_state == FeedSourceSubscribeStates.Subscribed;
      _subscribe_self = new RelayCommand<bool>(ExecuteSubscribeSelf);
    }
    public string Url { get { return _.uri; } }
    public string Title { get { return _.name; } }
    public string Summary { get { return _.description; } }
    public string Website { get { return _.website; } }

    public bool HasSubscribed { get { return _has_subscribed; } set { Set(ref _has_subscribed, value); } }
    public ICommand SubscribeCommand
    {
      get { return _subscribe_self ;}
    }

    async void ExecuteSubscribeSelf(bool sub)
    {
      if ((sub && HasSubscribed) || (!sub && !HasSubscribed))
        return;
      IsBusying = true;
      if (sub)
      {
        Debug.Assert(!string.IsNullOrEmpty(Url));
        var uri = BackendService.Compile(ServiceLocator.BackendAddress(), BackendService.FeedSourceSubscribe, new { uri = Url });

//        var rel = "/api/feed_source/subscribe.json?uri=" + Uri.EscapeDataString(Url);
        var v = await utils.HttpClientUtils.Get<FeedSource>(uri);
        if (v.code != 0)
        {
          Reason = v.reason;
          MessengerInstance.Send(new BackendError { code = v.code, reason = v.reason });
          return;
        }
        MessengerInstance.Send(new messages.SubscribeFeedSource { source = v.data });
        _.subscribe_state = FeedSourceSubscribeStates.Subscribed;
      }
      else
      {
        MessengerInstance.Send(new messages.UnsubscribeFeedSource { source = Url});
      }
      HasSubscribed = sub;
      IsBusying = false;
      /*	http.HandleFunc("/api/feed_source/unsubscribe.json", webapi_feedsource_unsubscribe)*/
    }
  }
}
