using GalaSoft.MvvmLight.Messaging;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Runtime.Serialization;
using System.Text;
using System.Threading.Tasks;
using System.Windows;

namespace famousfront.core
{
    [DataContract]
    class Flags
    {
        public Flags()
        {
            Backend = "127.0.0.1:8002";
            KaPeriod = 100;  // milliseconds
            ContentMarginHeight = 16.0;
            FeedSourceWidth = 360.0;
            ContentMarginMinWidth = 180.0;
        }
        [DataMember(Name="backend")]
        public string Backend
        {
            get;
            set;
        }
        [DataMember(Name = "ka_period")]
        public int KaPeriod
        {
            get;
            set;
        }
        [DataMember(Name = "contentmargin_height")]
        public double ContentMarginHeight { get; set; }

        [DataMember(Name = "feedsource_width")]
        public double FeedSourceWidth { get; set; }

        [DataMember(Name = "contentmargin_minwidth")]
        public double ContentMarginMinWidth { get; set; }
    }

}
