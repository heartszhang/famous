using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using famousfront.datamodels;
using System.Windows.Documents;
using System.Windows.Markup;
using System.Xml;
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
        /*
        FlowDocument _doc;
        public FlowDocument document
        {
            get { return _doc ?? (_doc = load_flowdoc()); }
        }
        FlowDocument load_flowdoc()
        {
            return XamlReader.Load(new XmlTextReader(new System.IO.StringReader(_.summary))) as FlowDocument;
        }
         * */
    }
}
