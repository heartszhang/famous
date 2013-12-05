using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class BackendError
  {
    [DataMember]
    internal int code { get; set; }
    [DataMember]
    internal string reason { get; set; }
  }
}