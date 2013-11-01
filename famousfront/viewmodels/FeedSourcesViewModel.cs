using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.viewmodels
{
    using FeedSources = ObservableCollection<FeedSourceViewModel>;
    class FeedSourcesViewModel : famousfront.core.ViewModelBase
    {
        FeedSources _sources = new FeedSources();
        public FeedSources Sources
        {
            get { return _sources; }
        }
    }
}
