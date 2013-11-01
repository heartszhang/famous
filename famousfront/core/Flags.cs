using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;

namespace famousfront.core
{
    class Flags
    {
        public Flags()
        {
        }
        public async void AsyncLoad()
        {
            await Task.Run(() => DoLoad());
        }
        public async void AsyncSave()
        {
            await Task.Run(() => DoSave());
        }
        private void DoSave()
        {

        }
        private void DoLoad()
        {

        }
    }

}
