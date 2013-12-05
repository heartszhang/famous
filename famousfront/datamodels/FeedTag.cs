using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedTag
  {
    [DataMember(Name = "value")]
    public string Name
    {
      get;set;
    }
  }
}