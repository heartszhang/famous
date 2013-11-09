using famousfront.core;
using famousfront.datamodels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.viewmodels
{
    class ImageUnitViewModel : ViewModelBase
    {
        FeedMedia _;
        internal ImageUnitViewModel(FeedMedia v)
        {
            _ = v;
        }
        public double Scale
        {
            get { return _.height > 0 ? (double)_.width / _.height : 0.0; }
        }
        public string Url
        {
            get { return _.thumbanil; }
        }
    }
}
