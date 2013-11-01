using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.viewmodels
{
    using FeedEntries = System.Collections.ObjectModel.ObservableCollection<FeedEntryViewModel>;
    class FeedEntriesViewModel : famousfront.core.ViewModelBase
    {
        FeedEntries _entries = new FeedEntries();
        public FeedEntries Entries { get { return _entries; } }
    }
}
