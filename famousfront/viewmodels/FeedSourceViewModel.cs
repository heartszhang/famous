using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using famousfront.datamodels;
namespace famousfront.viewmodels
{
    internal class FeedSourceViewModel : famousfront.core.ViewModelBase
    {
        FeedSource _ = new FeedSource();
        public string Category { get { return _category ?? first_or_default_category(); }  }
        public string Name { get { return _.name; } set { _.name = value; RaisePropertyChanged() ; } }
        public string Uri { get { return _.uri; }  }

        string _category = null;
        private string first_or_default_category()
        {
            return _.categories == null || _.categories.Length <= 0? "" : _.categories[0];
        }
        private bool append_category(string val)
        {
            if (val == first_or_default_category())
            {
                return false;
            }
            var n = new string[] { val };
            _.categories = n.Concat(_.categories).ToArray();
            return true;
        }
    }
}
