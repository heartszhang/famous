using famousfront.viewmodels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.messages
{
    internal class GenericMessage
    {
        internal int code;
        internal string reason = null;
        protected GenericMessage()
        {
            code = 0;
        }
    }
    internal class BackendInitialized : GenericMessage
    {
    }
    internal class BackendShutdown : GenericMessage
    {

    }
    internal class BackendInitializing : GenericMessage { }
    internal class BackendError : GenericMessage { }

    internal class VideoPlayRequest : GenericMessage
    {
        internal string source = null;
    }
    internal class ToggleFeedSource : GenericMessage
    {

    }
    internal class DropFeedSource : GenericMessage{
        internal FeedSourceViewModel model;
    }
    internal class SubscribeFeedSource : GenericMessage { internal datamodels.FeedSource source;}
    internal class UnsubscribeFeedSource : GenericMessage { internal string source;}
    internal class ShowFindFeedSourceView : GenericMessage
    {

    }
    internal class ImageTipRequest : GenericMessage
    {
      internal datamodels.FeedMedia image;
      internal bool open;
    }
}