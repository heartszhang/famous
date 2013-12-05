using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedTitle
  {
    [DataMember(EmitDefaultValue = false)]
    public string main
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string[] others
    {
      get;set;
    }
  }
}