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
    }

}
