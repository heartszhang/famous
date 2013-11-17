using famousfront.core;
using famousfront.datamodels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.viewmodels
{
  class FeedSourceFindEntryViewModel : TaskViewModel
  {
    FeedSourceFindEntry _;
    internal FeedSourceFindEntryViewModel(FeedSourceFindEntry v)
    {
      _ = v;
    }
    public string Uri { get { return _.url; } }
    public string Title { get { return _.title; } }
    public string Summary { get { return _.summary; } }
    public string Website { get { return _.website; } }
    public bool HasSubscribed { get { return _.subscribed; } }
  }
}
