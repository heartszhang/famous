using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using famousfront.datamodels;
using System.Windows.Documents;
using System.Windows.Markup;
using System.Xml;
using famousfront.core;
using System.Diagnostics;
namespace famousfront.viewmodels
{
    class FeedEntryViewModel : famousfront.core.ViewModelBase
    {
        static readonly System.DateTime utime = new DateTime(1970, 1, 1, 0, 0, 0, 0);
        FeedEntry _ = new FeedEntry();

        internal FeedEntryViewModel(FeedEntry v)
        {
            _ = v;
            _pub_day = publish_day();
            HasDocument = (_.status & FeedStatuses.Feed_status_text_empty) == 0;
            var inline = _.status & FeedStatuses.Feed_status_media_inline;
            var imgone = _.status & FeedStatuses.Feed_status_image_one;
            var imgmany = _.status & FeedStatuses.Feed_status_image_many;
            var media = _.status & (FeedStatuses.Feed_status_media_one | FeedStatuses.Feed_status_media_many);
            if (imgone != 0 && inline == 0 && media == 0) {
                Media = new ImageElementViewModel(_.images[0]);
            }
            else if (imgmany != 0 && inline == 0 && media == 0)
            {
                Media = new ImageGalleryViewModel(_.images);
            }
            if (media != 0)
            {
                Media = new MediaElementViewModel(_.videos[0], (imgone | imgmany) != 0 ? _.images[0] : null);
            }
        }
        public bool HasMedia
        {
          get { return Media != null; }
        }
        bool is_media_inline()
        {
          var inline = _.status & FeedStatuses.Feed_status_media_inline;
          return inline != 0;
        }
        public bool HasImageGallery
        {
          get { return  !is_media_inline() && (_.status & FeedStatuses.Feed_status_image_many) != 0; }
        }
        public bool HasVideo
        {
          get { return !is_media_inline() && (_.status & (FeedStatuses.Feed_status_media_one | FeedStatuses.Feed_status_media_many)) != 0; }
        }
      public bool HasImageOne
      {
        get { return (_.status & FeedStatuses.Feed_status_image_one) != 0; }
      }
        public string Summary { get { return _.summary; } }
        public string Title { get { return _.title.main; } }

        string _pub_day = null;
        public string PubDay { get { return _pub_day ; } }

        string publish_day()
        {
            var p = utime.AddMilliseconds(_.pubdate / 1000000);
            return p.ToString("D");
        }
        TaskViewModel _media;
        public TaskViewModel Media
        {
            get { return _media; }
            private set {Set(ref _media, value);}
        }
        bool _has_document = true;
        public bool HasDocument
        {
            get { return _has_document; }
            private set { Set(ref _has_document, value); }
        }
    }
}
