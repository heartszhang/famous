using famousfront.core;
using famousfront.datamodels;
using GalaSoft.MvvmLight.Command;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Input;

namespace famousfront.viewmodels
{
  class FeedSourceFindEntryViewModel : TaskViewModel
  {
    FeedSourceFindEntry _;
    internal FeedSourceFindEntryViewModel(FeedSourceFindEntry v)
    {
      _ = v;
    }
    public string Url { get { return _.url; } }
    public string Title { get { return _.title; } }
    public string Summary { get { return _.summary; } }
    public string Website { get { return _.website; } }
    public bool HasSubscribed { get { return _.subscribed; } set { _.subscribed = value; RaisePropertyChanged(); } }
    ICommand _subscribe_self;
    public ICommand SubscribeCommand
    {
      get { return _subscribe_self ?? (_subscribe_self = subscribe_self());}
    }
    ICommand subscribe_self()
    {
      return new RelayCommand<bool>(ExecuteSubscribeSelf);
    }
    async void ExecuteSubscribeSelf(bool sub)
    {
      if ((sub && HasSubscribed) || (!sub && !HasSubscribed))
        return;
      IsBusying = true;
      if (sub)
      {
        Debug.Assert(!string.IsNullOrEmpty(Url));
        var rel = "/api/feed_source/subscribe.json?uri=" + Uri.EscapeDataString(Url);
        var v = await famousfront.utils.HttpClientUtils.Get<FeedSource>(ServiceLocator.BackendPath(rel));
        if (v.code != 0)
        {
          Reason = v.reason;
          MessengerInstance.Send(new BackendError() { code = v.code, reason = v.reason });
          return;
        }
        MessengerInstance.Send(new messages.SubscribeFeedSource { source = v.data });
        HasSubscribed = true;
      }
      else
      {
        MessengerInstance.Send(new messages.UnsubscribeFeedSource { source = Url});
      }
      IsBusying = false;
      /*	http.HandleFunc("/api/feed_source/unsubscribe.json", webapi_feedsource_unsubscribe)*/
    }
  }
}
