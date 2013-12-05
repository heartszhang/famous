using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedAuthor
  {
    [DataMember(EmitDefaultValue = false)]
    public string name
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string email
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public ulong id
    {
      get;set;
    }
  }
}