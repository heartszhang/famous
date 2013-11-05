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
    }
}
