using GalaSoft.MvvmLight.Messaging;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.CompilerServices;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.core
{
    public abstract class ViewModelBase : GalaSoft.MvvmLight.ViewModelBase
    {
        protected ViewModelBase()
        {
        }

        protected ViewModelBase(IMessenger messenger)
            : base(messenger)
        {
        }

        protected new IMessenger MessengerInstance
        {
            get { return base.MessengerInstance ?? Messenger.Default; }
        }

        public override void Cleanup()
        {
            MessengerInstance.Unregister(this);            
        }

        protected override void RaisePropertyChanged([CallerMemberName] string propertyName = null)
        {
            base.RaisePropertyChanged(propertyName);
        }
        protected void Set<T>(ref T field, T value, [CallerMemberName] string name = null)
        {
            Set(name, ref field, value);
        }
    }
}
