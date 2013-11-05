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
        FeedSource _ = null;
        internal FeedSourceViewModel(FeedSource val)
        {
            _ = val;
            News = _.description;
        }
        public string Category { get { return _category ?? first_or_default_category(); }  }
        public string Name { get { return _.name; } set { _.name = value; RaisePropertyChanged() ; } }
        public string Uri { get { return _.uri; }  }
        public string Description { get { return _.description; } }

        string _logo = null;
        public string Logo { get { return _logo; } private set {  Set(ref _logo, value);  } }

        string _news;
        public string News { get { return _news; } private set {  Set(ref _news, value);  } }

        public int UnreadsCount 
        {
            get { return _.unreaded; } 
            private set { if (_.unreaded == value) return; _.unreaded = value; RaisePropertyChanged(); } 
        }

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
        private bool has_logo()
        {
            return false;
        }
    }
}
