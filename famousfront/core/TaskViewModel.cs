using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.core
{
    class TaskViewModel : ViewModelBase
    {
        bool _is_busying;
        public bool IsBusying
        {
            get { return _is_busying; }
            set { Set(ref _is_busying, value);  }
        }
        bool _is_ready;
        public bool IsReady
        {
            get { return _is_ready; }
            set { Set(ref _is_ready, value); }
        }

        string _reason;
        public string Reason
        {
            get { return _reason; }
            set { Set(ref _reason, value); }
        }
    }
}
